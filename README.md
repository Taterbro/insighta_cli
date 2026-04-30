# Insighta CLI

A command-line interface for managing user profiles on the Insighta Labs server. Built with Go and Cobra, this CLI provides seamless authentication, profile management, and data export capabilities.

**HNG Internship 14 - Stage 3 Backend Track Submission**

## Table of Contents

- [Features](#features)
- [Installation](#installation)
- [Quick Start](#quick-start)
- [Commands](#commands)
  - [Authentication](#authentication)
  - [Profile Management](#profile-management)
  - [User Information](#user-information)
- [Configuration](#configuration)
- [Examples](#examples)
- [Requirements](#requirements)
- [Contributing](#contributing)
- [License](#license)

## Features

✨ **Core Capabilities:**

- **GitHub OAuth Authentication** - Secure login via GitHub OAuth integration
- **Profile Management** - Create, view, and manage user profiles
- **Profile Export** - Export user profiles in JSON format
- **Session Management** - Persistent authentication with token-based sessions
- **Table Display** - Beautiful formatted output for profile listings
- **Auto-browser Launch** - Automatic browser opening for OAuth flows

## Installation

### Prerequisites

- **Go** 1.25.3 or higher ([Download Go](https://golang.org/dl/))
- **Git** for version control

### Build from Source

```bash
# Clone the repository
git clone https://github.com/Taterbro/insighta.git
cd insighta

# Install dependencies
go mod tidy

# Install globally
go install .
```

### Using the Binary

After building, you can use the CLI:

## Quick Start

### 1. Authenticate with GitHub

```bash
insighta login
```

This command will:

- Open your default browser to authenticate with GitHub
- Prompt you to authorize the Insighta Labs application
- Save your authentication tokens locally for future use

### 2. Verify Your Login

```bash
insighta whoami
```

Displays your current user information including username, email, and avatar.

### 3. Manage Profiles

```bash
# List all profiles
insighta profiles

# Create a new profile
insighta profiles create --name "John Doe"

# Export a profile
insighta profiles export --profile-id <ID> --format json

# Get a specific profile
insighta profiles get --profile-id <ID>
```

### 4. Logout

```bash
insighta logout
```

Clears your local session and authentication tokens.

## Commands

### Authentication

#### `insighta login`

Authenticate with the Insighta Labs server via GitHub OAuth.

**Usage:**

```bash
insighta login
```

**What it does:**

- Launches your browser for GitHub OAuth authentication
- Stores access and refresh tokens locally
- Saves your user details for quick reference

**Output:**

```
✓ Successfully logged in as username
```

---

#### `insighta logout`

Clear your authentication credentials and end your session.

**Usage:**

```bash
insighta logout
```

**Output:**

```
✓ Successfully logged out
```

---

#### `insighta whoami`

Display information about the currently authenticated user.

**Usage:**

```bash
insighta whoami
```
