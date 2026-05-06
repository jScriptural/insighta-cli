# insighta-cli

**The official Command Line Interface for Insighta** — a powerful profile discovery and intelligent query engine.

`insighta-cli` lets you search, filter, export, and manage profiles using both structured flags and natural language — all from your terminal.

---

## Features

- GitHub OAuth login (PKCE flow)
- Natural language search powered by the Intelligent Query Engine ([Intelligent Query Engine](https://github.com/jscriptural/intelligent-query-engine))
- Structured filtering (age, gender, country, etc.)
- Profile operations: list, search, get, create
- Data export to CSV or JSON
- Automatic token refresh and session management
- Clean, formatted JSON output

---

## Usage

### Login with GitHub
insighta login

### Check logged-in user
insighta whoami

### Logout
insighta logout


### List profiles with filters
insighta profiles list --min-age 18 --gender female --limit 20

### Natural language search
insighta profiles search "young males from Nigeria"

### Combined usage
insighta profiles search "young males" --min-age 18 --sort-by age --order ASC

### Get a specific profile
insighta profiles get <profile-id>

### Create a profile
insighta profiles create --name "John Doe"

# Export to CSV or JSON
``` bash
insighta export csv --gender male --min-age 25 --limit 15
insighta export json --country-name Nigeria
