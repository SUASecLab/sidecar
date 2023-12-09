dnl Tool for generating access rules for the SUASecLab
dnl
dnl Define rules
define(`RULE', `        {"key": "$1", "value": "$2"}$3')dnl
dnl
dnl Define permissions
define(`ALLOWED', `allowed')dnl
define(`MODERATOR', `moderator')dnl
define(`DENIED', `denied')dnl
dnl Generate file
{
    "rules": [
RULE(`bbbModerator', `MODERATOR', `,')
RULE(`jitsiModerator', `MODERATOR', `,')
RULE(`noVNC', `ALLOWED', `,')
RULE(`showComponents', `ALLOWED', `,')
RULE(`showVideo', `ALLOWED', `,')
RULE(`updateComponents', `MODERATOR', `,')
RULE(`updateVideo', `ALLOWED', `')
    ]
}
