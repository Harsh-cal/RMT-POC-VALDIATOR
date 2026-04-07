# RMT POC Validator

Release Management Tool (RMT) POC for validating aircraft software releases before deployment.

This system helps engineering teams answer one core question early:
Is this release safe to deploy to this fleet and aircraft context?

---

## Why This Exists

Modern aircraft software releases include many containers, dependencies, and compatibility constraints.

Without a structured validator, teams face:
- Manual and inconsistent checks
- Missed incompatibilities before deployment
- Late-stage rollback risk
- Limited visibility into release readiness

RMT solves this by running automated pre-deployment validation with clear risk signals.

---

## What The System Does

- Validates release payloads before deployment
- Detects version conflicts and missing dependencies
- Applies maturity-based risk checks
- Computes final status and risk level
- Generates recommendations for remediation
- Produces AI insight summary (with fallback when AI key is absent)
- Stores validation history in MongoDB
- Shows dashboard analytics (history, trends, recurring issues)
- Sends Telegram alerts for both pass and fail validations with fleet and targeting context
- Supports export reports in CSV, PDF, and Excel formats

---

## Key Value Delivered

- Faster go or no-go release decisions
- Reduced operational deployment risk
- Better traceability and auditability
- Less manual triage effort
- Higher confidence for fleet-specific rollouts

---

## Tech Stack

Backend
- Go
- Gin (REST API)
- MongoDB
- OpenAI API integration (gpt-4o-mini via SDK)
- gofpdf and excelize for exports

Frontend
- React
- Vite
- Tailwind CSS
- Axios
- Lucide React icons

Notifications
- Telegram Bot API

---

## High-Level Architecture

1. User submits release payload from UI or API
2. Backend runs validation engine
3. System calculates risk and status
4. Recommendations and insight are generated
5. Result is persisted in MongoDB
6. Telegram alert is triggered asynchronously
7. Frontend dashboard reads analytics endpoints and visualizes trends

---

## Installation From Scratch

## 1) Clone project

    git clone <your-repository-url>
    cd RMT-POC-Validator

## 2) Start MongoDB

Run local MongoDB, or point to your remote MongoDB instance.

Default expected URI:
- mongodb://localhost:27017

## 3) Backend setup

    cd backend
    go mod tidy

Create backend .env file with:

    PORT=8080
    MONGODB_URI=mongodb://localhost:27017
    DB_NAME=rmt_validator
    OPENAI_API_KEY=your_openai_key_optional
    TELEGRAM_BOT_TOKEN=your_telegram_bot_token_optional
    TELEGRAM_CHAT_ID=your_telegram_chat_id_optional

Run backend server:

    go run main.go

Backend will start on:
- http://localhost:8080

## 4) Frontend setup

Open a new terminal:

    cd frontend
    npm install
    npm run dev

Frontend will run on:
- http://localhost:5173

---

## API Overview

Base path:
- /api/dev/v1

Core endpoints:
- POST /validate
- POST /validate/chat
- POST /validate/export
- GET /releases/history
- GET /releases/trends
- GET /issues/recurring
- GET /releases/:id1/compare/:id2

---

## Analytics Behavior Notes

- History endpoint returns recent records with pagination
- Trends endpoint is day-wise over selected day range
- Recurring Issues is cumulative frequency over selected day range
  - Example: 70x means that issue appeared 70 times in that selected window

---

## Telegram Alert Behavior

If Telegram env variables are configured, each validation sends:
- Release metadata
- Status and risk
- Fleet and aircraft targeting details
- Full issue list (not truncated)

If Telegram config is missing, validation still works normally.

---

## AI Insight Behavior

If OPENAI_API_KEY is configured:
- AI-generated release insight is returned

If not configured or API fails:
- Deterministic fallback insight is returned
- System continues to function without blocking validation

---

## Suggested Demo Flow

1. Open frontend
2. Validate a mock release
3. Observe risk, issues, and recommendations
4. Check Telegram alert
5. Open Release History dashboard
6. Review trends and recurring issues
7. Export validation report

---

## Troubleshooting

Backend says MongoDB not initialized
- Verify MONGODB_URI and MongoDB service availability

Port already in use
- Change PORT in .env and restart backend

No Telegram messages
- Verify TELEGRAM_BOT_TOKEN and TELEGRAM_CHAT_ID

AI insight missing
- Verify OPENAI_API_KEY
- Fallback insight will still be shown when key is absent

---

## Future Enhancements

- Better fix-rate logic based on true issue resolution lifecycle
- Advanced filtering by fleet, aircraft type, and system
- Role-based access and audit trails
- CI pipeline integration for release gating

---

## A quick demo 

![chrome-capture-2026-04-06 (1)](https://github.com/user-attachments/assets/9fa88ea0-fdbf-421b-b9a8-0af89cefe167)
