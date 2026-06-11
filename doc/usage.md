# Usage

Navigate to any Git repository and run:

```bash
gocommit
```

You'll be prompted to enter a commit message. The tool executes:

1. `git add .`
2. `git commit -m "your message"`
3. `git push`

## Example

```bash
$ gocommit
Enter commit message: fix: correct typo in README
[main b2n3m45n] fix: correct typo in README
 1 file changed, 1 insertion(+), 1 deletion(+)
To github.com:your-user/your-repo.git
```