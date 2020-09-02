你好！
很冒昧用这样的方式来和你沟通，如有打扰请忽略我的提交哈。我是光年实验室（gnlab.com）的HR，在招Golang开发工程师，我们是一个技术型团队，技术氛围非常好。全职和兼职都可以，不过最好是全职，工作地点杭州。
我们公司是做流量增长的，Golang负责开发SAAS平台的应用，我们做的很多应用是全新的，工作非常有挑战也很有意思，是国内很多大厂的顾问。
如果有兴趣的话加我微信：13515810775  ，也可以访问 https://gnlab.com/，联系客服转发给HR。
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
