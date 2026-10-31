#!/usr/bin/env bash
# Behavioural tests for 000001_init.up.sql (62 checks).
# Usage: apply the migration to an EMPTY database, then: DB=<database> bash db/test_constraints.sh
# It inserts test rows, so never run it against real data.
DB=${DB:-bl}; pass=0; fail=0
q(){ mysql -N "$DB" -e "$1" 2>&1; }
expect_ok(){ out=$(q "$2"); if echo "$out" | grep -q "^ERROR"; then echo "FAIL  $1 -> $out"; fail=$((fail+1)); else echo "PASS  $1"; pass=$((pass+1)); fi; }
expect_err(){ out=$(q "$3"); if echo "$out" | grep -q "ERROR $2"; then echo "PASS  $1 (ERROR $2)"; pass=$((pass+1)); else echo "FAIL  $1 expected ERROR $2 got: $out"; fail=$((fail+1)); fi; }
expect_val(){ out=$(q "$3"); if [ "$out" = "$2" ]; then echo "PASS  $1 = $2"; pass=$((pass+1)); else echo "FAIL  $1 expected '$2' got '$out'"; fail=$((fail+1)); fi; }

# --- seed minimal data ---
expect_ok "seed users"     "INSERT INTO users (id,email,password_hash,role) VALUES (1,'admin@brightloop.test','x','admin'),(2,'staff@brightloop.test','x','staff'),(3,'Cust@Shop.test','x','customer');"
expect_ok "seed customer"  "INSERT INTO customers (id,user_id,name) VALUES (1,3,'Small Shop'),(2,NULL,'Big Corp');"
expect_ok "seed settings"  "INSERT INTO company_settings (id) VALUES (1);"
expect_ok "seed catalog"   "INSERT INTO products (id,name,product_type,sales_price,cost_price) VALUES (1,'CRM Suite','software',1000,400); INSERT INTO plans (id,name,price,billing_period,available_from) VALUES (1,'Monthly','500.00','monthly','2026-01-01'); INSERT INTO payment_terms (id,name,due_days) VALUES (1,'Net 15',15); INSERT INTO close_reasons (id,name,is_churn) VALUES (1,'Renewed',FALSE),(2,'Too expensive',TRUE);"
expect_ok "seed taxes"     "INSERT INTO tax_rules (id,name,rate_percent,is_default) VALUES (1,'GST 18%',18,TRUE),(2,'Cess 1%',1,FALSE);"
expect_ok "seed discount"  "INSERT INTO discount_rules (id,name,discount_type,value,scope,starts_on,ends_on,max_uses,created_by) VALUES (1,'Seats 10%','percentage',10,'subscription','2026-01-01','2026-12-31',2,1);"
expect_ok "seed sub"       "INSERT INTO subscriptions (id,number,customer_id,plan_id,payment_terms_id,status,start_date) VALUES (1,'SUB-2026-00001',1,1,1,'active','2026-03-01');"

# --- A: identity ---
expect_err "email unique ignores case"            1062 "INSERT INTO users (email,password_hash,role) VALUES ('cust@shop.TEST','x','customer');"
expect_err "company_settings is single-row"       3819 "INSERT INTO company_settings (id) VALUES (2);"
# --- B: catalog ---
expect_err "only one default tax rule"            1062 "INSERT INTO tax_rules (name,rate_percent,is_default) VALUES ('VAT 5%',5,TRUE);"
expect_ok  "many non-default tax rules allowed"        "INSERT INTO tax_rules (name,rate_percent,is_default) VALUES ('VAT 12%',12,FALSE),('VAT 0%',0,FALSE);"
expect_err "tax rate over 100"                    3819 "INSERT INTO tax_rules (name,rate_percent) VALUES ('Bad',101);"
expect_err "percentage discount over 100"         3819 "INSERT INTO discount_rules (name,discount_type,value,scope,starts_on,ends_on,created_by) VALUES ('Bad','percentage',150,'subscription','2026-01-01','2026-02-01',1);"
expect_err "discount end before start"            3819 "INSERT INTO discount_rules (name,discount_type,value,scope,starts_on,ends_on,created_by) VALUES ('Bad','fixed',10,'subscription','2026-02-01','2026-01-01',1);"
expect_err "plan availability reversed"           3819 "INSERT INTO plans (name,price,billing_period,available_from,available_until) VALUES ('Bad',1,'weekly','2026-05-01','2026-04-01');"
expect_err "negative sales price"                 3819 "INSERT INTO products (name,product_type,sales_price,cost_price) VALUES ('Bad','software',-1,0);"
expect_ok  "discount used up to max_uses"              "UPDATE discount_rules SET times_used = 2 WHERE id = 1;"
expect_err "discount cannot exceed max_uses"      3819 "UPDATE discount_rules SET times_used = 3 WHERE id = 1;"
expect_err "times_used cannot go below zero"      1690 "UPDATE discount_rules SET times_used = times_used - 3 WHERE id = 1;"
# --- C: subscriptions ---
expect_err "expiry before start"                  3819 "INSERT INTO subscriptions (number,customer_id,plan_id,payment_terms_id,start_date,expiry_date) VALUES ('SUB-X1',1,1,1,'2026-03-01','2026-02-01');"
expect_err "closed requires closed_at"            3819 "UPDATE subscriptions SET status='closed' WHERE id=1;"
expect_err "duplicate subscription number"        1062 "INSERT INTO subscriptions (number,customer_id,plan_id,payment_terms_id,start_date) VALUES ('SUB-2026-00001',1,1,1,'2026-03-01');"
expect_ok  "first renewal quotation"                   "INSERT INTO subscriptions (id,number,customer_id,plan_id,payment_terms_id,status,start_date,renewed_from_id) VALUES (2,'SUB-2026-00002',1,1,1,'quotation','2027-03-01',1);"
expect_err "second open renewal blocked"          1062 "INSERT INTO subscriptions (number,customer_id,plan_id,payment_terms_id,status,start_date,renewed_from_id) VALUES ('SUB-2026-00003',1,1,1,'quotation','2027-03-01',1);"
expect_ok  "close first renewal"                       "UPDATE subscriptions SET status='closed', closed_at=NOW(3), close_reason_id=2 WHERE id=2;"
expect_ok  "new renewal allowed after close"           "INSERT INTO subscriptions (id,number,customer_id,plan_id,payment_terms_id,status,start_date,renewed_from_id) VALUES (3,'SUB-2026-00003',1,1,1,'quotation','2027-03-01',1);"
expect_ok  "subscription line"                         "INSERT INTO subscription_lines (subscription_id,product_id,quantity,unit_price) VALUES (1,1,3,'33.33');"
expect_val "line amount generated"           "99.99"   "SELECT amount FROM subscription_lines WHERE subscription_id=1;"
expect_err "zero quantity line"                   3819 "INSERT INTO subscription_lines (subscription_id,product_id,quantity,unit_price) VALUES (1,1,0,10);"
# --- D: invoices ---
expect_ok  "draft invoice without number"              "INSERT INTO invoices (id,customer_id,subscription_id,invoice_date,due_date,billing_name,subtotal,tax_total,total) VALUES (1,1,1,'2026-03-01','2026-03-16','Small Shop',1000,180,1180);"
expect_ok  "second draft without number (NULLs not unique)" "INSERT INTO invoices (id,customer_id,invoice_date,due_date,billing_name,total) VALUES (2,1,'2026-03-01','2026-03-16','Small Shop',0);"
expect_err "confirmed invoice needs a number"     3819 "UPDATE invoices SET status='confirmed' WHERE id=1;"
expect_ok  "confirm with number"                       "UPDATE invoices SET status='confirmed', number='INV-2026-00001', confirmed_at=NOW(3) WHERE id=1;"
expect_err "duplicate invoice number"             1062 "UPDATE invoices SET status='confirmed', number='INV-2026-00001' WHERE id=2;"
expect_err "due date before invoice date"         3819 "UPDATE invoices SET due_date='2026-02-01' WHERE id=1;"
expect_val "payment_status not_paid"     "not_paid"     "SELECT payment_status FROM invoices WHERE id=1;"
expect_ok  "partial payment"                           "UPDATE invoices SET amount_paid=500 WHERE id=1;"
expect_val "payment_status partial"      "partial"      "SELECT payment_status FROM invoices WHERE id=1;"
expect_err "paid status requires full payment"    3819 "UPDATE invoices SET status='paid' WHERE id=1;"
expect_err "overpayment blocked"                  3819 "UPDATE invoices SET amount_paid=1180.01 WHERE id=1;"
expect_ok  "full payment then paid"                    "UPDATE invoices SET amount_paid=1180, status='paid' WHERE id=1;"
expect_val "payment_status paid"         "paid"         "SELECT payment_status FROM invoices WHERE id=1;"
# invoice items
expect_ok  "recurring item"                            "INSERT INTO invoice_items (id,invoice_id,item_type,product_id,description,quantity,unit_price,amount) VALUES (1,1,'recurring',1,'CRM Suite',1,1000,1000);"
expect_ok  "discount item linked"                      "INSERT INTO invoice_items (id,invoice_id,item_type,linked_item_id,discount_rule_id,description,rate_percent,amount) VALUES (2,1,'discount',1,1,'Seats 10%',10,-100);"
expect_ok  "tax item linked"                           "INSERT INTO invoice_items (id,invoice_id,item_type,linked_item_id,tax_rule_id,description,rate_percent,amount) VALUES (3,1,'tax',1,1,'GST 18%',18,162);"
expect_err "discount item must be negative"       3819 "INSERT INTO invoice_items (invoice_id,item_type,linked_item_id,discount_rule_id,description,amount) VALUES (1,'discount',1,1,'Bad',5);"
expect_err "tax item must be linked"              3819 "INSERT INTO invoice_items (invoice_id,item_type,tax_rule_id,description,amount) VALUES (1,'tax',1,'Bad',5);"
expect_err "recurring item cannot be negative"    3819 "INSERT INTO invoice_items (invoice_id,item_type,description,amount) VALUES (1,'recurring','Bad',-5);"
expect_err "recurring item cannot carry tax rule" 3819 "INSERT INTO invoice_items (invoice_id,item_type,tax_rule_id,description,amount) VALUES (1,'recurring',1,'Bad',5);"
# billing guard, payments, discounts
expect_ok  "billing record"                            "INSERT INTO billing_records (subscription_id,period_start,period_end,invoice_id) VALUES (1,'2026-03-01','2026-03-31',1);"
expect_err "same period billed twice"             1062 "INSERT INTO billing_records (subscription_id,period_start,period_end,invoice_id) VALUES (1,'2026-03-01','2026-03-31',2);"
expect_ok  "payment recorded"                          "INSERT INTO payments (invoice_id,method,amount,payment_date,idempotency_key,recorded_by) VALUES (1,'upi',500,'2026-03-05','3f2b8c1e-0000-4000-8000-000000000001',2);"
expect_err "double submit same idempotency key"   1062 "INSERT INTO payments (invoice_id,method,amount,payment_date,idempotency_key,recorded_by) VALUES (1,'upi',500,'2026-03-05','3f2b8c1e-0000-4000-8000-000000000001',2);"
expect_err "zero payment"                         3819 "INSERT INTO payments (invoice_id,method,amount,payment_date,idempotency_key,recorded_by) VALUES (1,'cash',0,'2026-03-05','3f2b8c1e-0000-4000-8000-000000000002',2);"
expect_ok  "discount application"                      "INSERT INTO discount_applications (discount_rule_id,invoice_id,customer_id,subscription_id,amount) VALUES (1,1,1,1,100);"
expect_err "discount applied twice on one invoice" 1062 "INSERT INTO discount_applications (discount_rule_id,invoice_id,customer_id,subscription_id,amount) VALUES (1,1,1,1,100);"
# platform
expect_ok  "job with dedupe key"                       "INSERT INTO job_queue (queue_name,dedupe_key,payload,effective_at) VALUES ('email','invoice-sent:1','{\"invoice_id\":1}',NOW(3));"
expect_err "same job enqueued twice"              1062 "INSERT INTO job_queue (queue_name,dedupe_key,payload,effective_at) VALUES ('email','invoice-sent:1','{\"invoice_id\":1}',NOW(3));"
expect_ok  "jobs without dedupe key repeat freely"     "INSERT INTO job_queue (queue_name,payload,effective_at) VALUES ('email','{}',NOW(3)),('email','{}',NOW(3));"
# gapless counter: a rolled-back increment is not consumed
expect_ok  "counter init"                              "INSERT INTO number_sequences (name,next_value) VALUES ('invoice-2026',2);"
expect_ok  "counter rollback"                          "START TRANSACTION; SELECT next_value FROM number_sequences WHERE name='invoice-2026' FOR UPDATE; UPDATE number_sequences SET next_value=next_value+1 WHERE name='invoice-2026'; ROLLBACK;"
expect_val "counter unchanged after rollback" "2"      "SELECT next_value FROM number_sequences WHERE name='invoice-2026';"
echo "RESULT: $pass passed, $fail failed"
