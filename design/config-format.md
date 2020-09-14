# Config file format
## change log
- 10/07/2020 create first draft
- 21/08/2020 update config yaml format

## Proposal
- File format is `yaml` and file extension is `.yaml`
- Config file name is `config`
- The configuration file is placed in the below location
  - `${HOME}/.confg/sdbi/config.yaml`
- Format is below
  ``` yaml
  config:
      ConfigA:
          api-url: https://example.com
          user-token: user_token_xxxx
          ui-url: https://ui-example.com
      ConfigB:
          api-url: https://example.com
          user-token: user_token_zzzz
          ui-url: https://ui-example.com
  current: ConfigA  
  ```
  - "ConfigA" and "ConfigB" are the names of the user's Config.
    - Cannot include spaces in the config name.
  - user.token is personal token of Screwdriver.cd
  - The current is the name of the configuration you are currently using.