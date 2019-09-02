# atcoder-review-notifier  
Tool for reviewing problems solved with AtCoder on a regular basis.  
Notify solved problems 3 to 2 weeks ago by Slack.  
  
## Abstract  
- Notify review of AtCoder Problems user's soleved.  
    - Run periodically using Cloud Scheduler.  
- Architecture  
    - Cloud Scheduler --> Pub/Sub --> CloudFunction  
  
## Usage  
  
### Deploy Cloud Function  
  
```bash
gcloud functions deploy NotifyReview --project <gcp-project> \
  --runtime go111 \
  --trigger-topic <topic> \
  --set-env-vars ATCODER_USER=<AtCocer-user_id> \
  --set-env-vars SLACK_API_TOKEN=<slack-api-token> \
  --set-env-vars SLACK_CHANNEL=<slack-channel>
```
  
### Create Cloud Scheduler  
Please create Cloud Scheduler on console.  
- frequency: unix-cron format [(useful)](https://crontab.guru/)  
- target   : Pub/Sub  
- topic    : topic specified when Cloud Function deployed  
  
## Test Usage  
- test/
    ```bash
    $ go run main.go <AtCoder_USERNAME> <SLACK_API_TOKEN> <SLACK_CHANNEL>
    ```
