# Astropath
CLI to orchestrate the execution of local AI agents


# Intended Scope

- Different "roles" for Agents: Analyst, Developer, Tester, Reviewer, Debugger, Explorer. Each one of them is defined by a specialized prompt.
- Agents are coordinated through Markdown files, this adds context from each step and enables "human-in-the-loop".
- Leverage Claude's *-p* option, so agents actually finish after running.
- Agents have permissions to create git branches and modify code there, but not main.
- CLI with different commands to run each of the steps separately or as part of a bigger multiple-step task.

# To-Fix list

- Improve templates management
  - More 'dynamic' content in templates, like branches naming

- Improve the prompts so Claude:
  - Thinks less in useless stuff
  - Writes less words but more important ones

- Add some sort of logging to Claude usage (tokens, costs, queries, whatever)
- Improve user feedback while Claude is doing stuff. add support for stream output and redirect it to stdout and/or logs.
- Improve golang stuff
  - Reorganize cmd and Cobra handlers
  - I wanted to use channels for communication but agents execution is going to be sync'ed so may not needed

- **Idea**: different ASTROPATH.md files to track different tasks / executions (or attach them to the branches somehow to implement some tracing)

# Created By

- Fynardo
- Claude itself
