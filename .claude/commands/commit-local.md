Commit all staged and unstaged changes to git, but DO NOT push to remote.

Steps:
1. Run `git status` to see changes
2. Run `git diff` to understand what changed
3. Run `git log --oneline -3` to match commit message style
4. Stage relevant files (prefer specific files over `git add -A`)
5. Commit with a descriptive message following the repo's convention
6. Verify with `git status` that commit succeeded and branch is ahead of origin
7. DO NOT run `git push`

End by confirming: "Committed locally (not pushed): <short commit hash> <commit subject>"
