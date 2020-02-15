# Contribution Guide

Make sure to run linting & tests before pushed your commits. Follow these steps:

- Create `.env` containing:

  ```ini
  URL=https://sandbox.bca.co.id

  CLIENT_ID=xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx
  CLIENT_SECRET=xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx

  API_KEY=xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx
  API_SECRET=xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx

  CORPORATE_ID=BCAAPI2016
  ORIGIN_HOST=localhost

  CHANNEL_ID=95051
  CREDENTIAL_ID=BCAAPI

  LOG_PATH=bca.log
  ```

- Run linting & tests:

  ```shell
  export $(cat .env | xargs) && make ci
  ```
