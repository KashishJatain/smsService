#!/usr/bin/env bash
# ─────────────────────────────────────────────────────────────────────────────
# demo.sh  endd-to-end demonstration of the Polyglot SMS Service
# ─────────────────────────────────────────────────────────────────────────────
set -euo pipefail

SMS_SENDER="http://localhost:8080"
SMS_STORE="http://localhost:8081"
USER_ID="demo-user-$(date +%s)"

GREEN="\033[0;32m"; YELLOW="\033[1;33m"; CYAN="\033[0;36m"; RESET="\033[0m"

step() { echo -e "\n${CYAN}══ $* ${RESET}"; }
ok()   { echo -e "${GREEN}✔  $*${RESET}"; }
info() { echo -e "${YELLOW}ℹ  $*${RESET}"; }

step "1. Health checks"
curl -sf "$SMS_SENDER/v1/sms/send" -o /dev/null -w "" || true   # will 405, just to check reachability
curl -sf "$SMS_STORE/health" | python3 -m json.tool
ok "sms-store is healthy"

step "2. Send a successful SMS (userId=$USER_ID)"
RESPONSE=$(curl -sf -X POST "$SMS_SENDER/v1/sms/send" \
  -H "Content-Type: application/json" \
  -d "{\"userId\":\"$USER_ID\",\"phoneNumber\":\"+919876543210\",\"message\":\"Hello from demo!\"}")
echo "$RESPONSE" | python3 -m json.tool
MESSAGE_ID=$(echo "$RESPONSE" | python3 -c "import sys,json; print(json.load(sys.stdin)['messageId'])" 2>/dev/null || echo "N/A")
ok "SMS sent. messageId=$MESSAGE_ID"

step "3. Block the user"
curl -sf -X POST "$SMS_SENDER/v1/sms/block/$USER_ID" | python3 -m json.tool
ok "User blocked"

step "4. Attempt to send another SMS (should be BLOCKED)"
curl -sf -X POST "$SMS_SENDER/v1/sms/send" \
  -H "Content-Type: application/json" \
  -d "{\"userId\":\"$USER_ID\",\"phoneNumber\":\"+919876543210\",\"message\":\"This should be blocked\"}" \
  | python3 -m json.tool
ok "Blocked response received"

step "5. Unblock the user"
curl -sf -X DELETE "$SMS_SENDER/v1/sms/block/$USER_ID" | python3 -m json.tool
ok "User unblocked"

step "6. Send another SMS after unblocking"
curl -sf -X POST "$SMS_SENDER/v1/sms/send" \
  -H "Content-Type: application/json" \
  -d "{\"userId\":\"$USER_ID\",\"phoneNumber\":\"+919876543210\",\"message\":\"Back in action!\"}" \
  | python3 -m json.tool
ok "Second SMS sent"

step "7. Wait for Kafka consumer to persist events..."
sleep 3

step "8. Retrieve SMS history from Go SMS Store"
curl -sf "$SMS_STORE/v1/user/$USER_ID/messages" | python3 -m json.tool
ok "History retrieved"

echo -e "\n${GREEN}════════════════════════════════════════${RESET}"
echo -e "${GREEN}  End-to-end demo completed successfully!${RESET}"
echo -e "${GREEN}════════════════════════════════════════${RESET}\n"