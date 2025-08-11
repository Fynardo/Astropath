# Astropath
CLI to orchestrate the execution of local AI agents


# Intended Scope

- Different "roles" for Agents: Analyst, Developer, Tester, Reviewer, Debugger, Explorer. Each one of them is defined by a specialized prompt.
- Agents are coordinated through Markdown files, this adds context from each step and enables "human-in-the-loop".
- Leverage Claude's *-p* option, so agents actually finish after running.
- Agents have permissions to create git branches and modify code there, but not main.
- CLI with different commands to run each of the steps separately or as part of a bigger multiple-step task.

# To-Do list

- Automatically create ASTROPATH.md file if not exists (maybe ask to the user before doing it). This way, Claude only needs edit access.
- Check if .claude directory exists in the project


# Created By

- Fynardo
- Claude
