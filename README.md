# Todo API — CI/CD Learning Project

A simple Go Todo API project created mainly to learn and practice CI/CD pipeline.

## Purpose

This project focuses on understanding:

- GitHub Actions CI pipeline
- Automated testing
- Linting
- Docker image build
- Auto deployment using Render

## Tech Stack

- Go
- Docker
- GitHub Actions
- Render

## CI/CD Flow

```text
Push to main
   ↓
Run tests
   ↓
Run lint checks
   ↓
Build Docker image
   ↓
Render auto-deploys latest code



```md
## Live URL
https://todo-api-nozn.onrender.com

## Health Check
```bash
curl https://todo-api-nozn.onrender.com/health


```md
## Expected Response
```json
{
  "status": "my app is live"
}


```md
## Note
This project is not focused on business logic or production features.  
It was built mainly to practice CI/CD, Docker, GitHub Actions, and deployment flow.