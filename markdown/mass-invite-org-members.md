---
slug: 'orgapi_invitemember'
---

## How to invite several users to an organization

If you need to invite more than one user to an organization, you can definitely leverage the `InviteMember` API. To do that, a simple bash script would do the job by calling that endpoint for each invitee.

First, you'll create a bash script which will accept some parameters. Please note, this endpoint requires three attributes to be able to properly construct the payload: orgId, the user role and the user email, which means we need to pass them at runtime. So let's break the whole down to three simple steps:

### Create a bash script

```bash
#! /bin/bash

orgId="$1"
arr="$2"
for email in "${arr[@]}";
  do
    curl -X POST -d '{ "orgId": "$orgId", "role": "Admin", "username": "$email" }' https://portal.percona.com/v1/orgs/"$orgId"/members
  done
```

You'll need to save it to a file, say "bulk-invite.sh" and make it executable:

```
chmod +x bulk-invite.sh
```

### Define the parameters to be passed to the script

invitees=("john.doe@acme.com" "phillip.myers@acme.com" "jessica.whales@acme.com")
orgId="a9c55b07-9962-490e-950a-7c4fb51081c9"

### Run the script

inviteMembers $orgId $invitees
