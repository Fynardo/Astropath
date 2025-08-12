# Astropath
CLI to orchestrate the execution of local AI agents


# Intended Scope

- Different "roles" for Agents: Analyst, Developer, Tester, Reviewer, Debugger, Explorer. Each one of them is defined by a specialized prompt.
- Agents are coordinated through Markdown files, this adds context from each step and enables "human-in-the-loop".
- Leverage Claude's *-p* option, so agents actually finish after running.
- Agents have permissions to create git branches and modify code there, but not main.
- CLI with different commands to run each of the steps separately or as part of a bigger multiple-step task.

# To-Do list

- Implement an 'astropath init' command that:
  - Automatically creates ASTROPATH.md file if not exists
  - Creates .claude directory if it doesn't exists
  - Populates the .claude/settings.json with 'Edit' and 'Git' perms

- Improve templates management
  - Move them to files for easier management and such
  - More 'dynamic' content in templates, like branches naming

- Improve the prompts so Claude:
  - Thinks less in useless stuff
  - Writes less words but more important ones
  - More clear directions of how to do things

- Add some sort of logging to tokens usage


# Created By

- Fynardo
- Claude itself
