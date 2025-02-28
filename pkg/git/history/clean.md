# clean commit message history
- list all hashes : `git log --pretty=format:%H`
- get commit information : `git show --no-color <commit-hash>`
- rewrite commit message and keep tags : `git filter-branch -f --msg-filter 'sed "s/release\/Version-[0-9].[0-9].[0-9]/develop/g"' --tag-name-filter cat -- --all`

```
rm -f /tmp/git;
touch /tmp/git;
git filter-branch \
    --subdirectory-filter <DIRECTORY> \
    --tag-name-filter cat \
    --commit-filter 'echo -n "s/${GIT_COMMIT}/" >>/tmp/git; \
                     NEW=`git_commit_non_empty_tree "$@"`; \
                     echo "${NEW}/g" >> /tmp/git; echo ${NEW}' \
    --msg-filter 'sed -f /tmp/git' \
    -- --all
```