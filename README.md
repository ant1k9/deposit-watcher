## Deposit Watcher

### Usage

```bash
>>> make
>>> ./deposit-migrate         # migration script to local sqlite database
>>> ./deposit-update          # get all updates for deposits
>>> ./deposit-app             # run tui application
>>> ./deposit-backup          # backup database and send zipped archive to dropbox
>>> ./deposit-rotate-backups  # rotate backups to keep only 30 last files
>>> make clean                # delete all executables deposit-*
```

### App commands

|Button             |Action                                             |
|-------------------|---------------------------------------------------|
|**\<Enter\>**      | open deposit detail page in your browser          |
|**\<r\>**          | reverse the sorting by rate                       |
|**\<d\>**          | remove deposit from the application permanently   |
|**\<PageDown\>**   | next page                                         |
|**\<PageUp\>**     | previous page                                     |
|**\<q\> \<C-c\>**  | exit from app                                     |

### Backup

It's needed to have `DB_BACKUP_TOKEN` in environment variables to send backup to dropbox.
