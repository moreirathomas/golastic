# Contributing guidelines

## Code formatting

- For go code, simply use go conventions (can be formatted via `go fmt`)
- Avoid _when possible_ lines of 80+ chars
- Every file should end with an empty line

## Git

- Branch names must follow this pattern: `<type>/name-of-related-issue`
- Commit messages must be 50 chars max (add a description if necessary)
- Commit messages must start with a verb in the present simple, such as `Add ...`, `Implement ...`, `Fix ...`

## Pull requests & Merge policy

- Do not create a Pull request if it is not directly related to an existing issue. Create an issue first if a change is needed.
- Merge must be performed using the method `squash`
- A Pull request must have been tested and reviewed before merging:
  - All tests must succeed
  - **2** reviewers at least must have approved the changes
