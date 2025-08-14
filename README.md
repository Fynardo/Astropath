# Astropath

**Astropath** is a CLI tool designed to orchestrate local AI agents for software development workflows, enabling human-in-the-loop coordination through specialized agent roles and markdown-based context sharing.

## Features & Capabilities

**Specialized Agent Roles**: Choose from different agent types, each optimized for specific development tasks:
- **Analyst**: Analyzes code structure, identifies patterns, and provides insights
- **Developer**: Implements features, fixes bugs, and writes code 
- **Explorer**: Navigates and documents codebase structure and functionality
- **Reviewer**: Performs code reviews and suggests improvements
- **Tester**: Creates and runs tests for validation

**Human-in-the-Loop Workflow**: Agents coordinate through `ASTROPATH.md` files that maintain context between execution steps, allowing you to guide and review the process at each stage.

**Git-Safe Operations**: Agents can create git branches and modify code in feature branches, but never directly modify the main branch, ensuring your codebase remains protected.

**Minimal Dependencies**: Built with Go standard library and Cobra CLI framework, keeping the tool lightweight and focused.

**Command-Based Interface**: Execute agents individually or as part of multi-step workflows using dedicated commands for each agent type.

## Usage Examples

### Initialize a Project
```bash
# Set up Astropath in your project directory
astropath init
```

### Explore Your Codebase
```bash
# Have the Explorer agent analyze and document your project structure
astropath explore "Please analyze the project structure and main components"
```

### Analyze Specific Code
```bash
# Use the Analyst to examine code patterns and potential issues
astropath analyze "Review the authentication module for security best practices"
```

### Implement New Features
```bash
# Let the Developer agent implement functionality based on requirements
astropath develop "Add user registration endpoint with validation"
```

### Code Review Process
```bash
# Get comprehensive code review feedback
astropath review "Please review the changes in the current branch"
```

### Multi-Step Workflow
```bash
# Use pipeline for coordinated multi-agent execution
astropath pipeline "Analyze, develop user authentication, then review the implementation"
```

### Raw Claude Interaction
```bash
# Direct interaction with Claude for custom tasks
astropath raw "Help me debug this specific function"
```

Each command creates or updates the `ASTROPATH.md` file with context, allowing you to review progress and provide guidance between steps.

# To-Fix list

- Improve the prompts so Claude:
  - Thinks less in useless stuff
  - Writes less words but more important ones

- Improve role branches management
  - If branch is specified -> tell the agent to use that
  - If no branch is specified but current branch is not main -> use that one
  - If no branch is specified and current branch is main -> Checkout to a new one

- Add some sort of logging to Claude usage (tokens, costs, queries, whatever)

- **Idea**: different ASTROPATH.md files to track different tasks / executions (or attach them to the branches somehow to implement some tracing)

# Created By

- Fynardo
- Claude itself
