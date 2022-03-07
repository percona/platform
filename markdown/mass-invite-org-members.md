---
slug: 'orgapi_invitemember'
---

## How to invite several users to an organization

If you need to invite more than one user to an organization, you can definitely leverage the `InviteMember` API. To do that, a simple bash script would do the job by calling that endpoint for each invitee.

First, you'll create a bash script which will accept some parameters. Please note, this endpoint requires three attributes to be able to properly construct the payload: orgId, the user role, the user email and the access token, which we need to pass at runtime. So let's break
this down to three simple steps:

### Create a bash script

```bash
#! /bin/bash

ORG_ID="$1"
MEMBERS="$2"
TOKEN="$3"

for EMAIL in "${MEMBERS[@]}";
  do
    curl -X POST -H "Content-Type: application/json" -H "Authorization: Bearer $TOKEN" -d '{ "orgId": "$ORG_ID", "role": "Admin", "username": "$EMAIL" }' https://portal.percona.com/v1/orgs/"$ORG_ID"/members
  done
```

You'll need to save it to a file, say "invite-members.sh" and make it executable:

```
chmod +x invite-members.sh
```

### Define the parameters to be passed to the script

orgId="a9c55b07-9962-490e-950a-7c4fb51081c9"
invitees=("john.doe@acme.com" "phillip.myers@acme.com" "jessica.whales@acme.com")
token="<ACCESS_TOKEN>"

Please note that you can get your acceess token by logging in to Percona Platform Portal and extracting the access token by running
the following chunk of code in the browser Dev Tools' `Console` tab:

```javascript
var oktaStorage = localStorage.getItem('okta-token-storage');
JSON.parse(oktaStorage).accessToken.accessToken;
```

### Run the script

./invite-members.sh $orgId $invitees $token
