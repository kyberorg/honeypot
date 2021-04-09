## More Information About This Repo
Boilerplated from [ MartinHeinz/go-project-blueprint ](https://github.com/MartinHeinz/go-project-blueprint)

### Adding New Libraries/Dependencies
```bash
go mod vendor
```

### Using GitHub Registry

Create and Push:

```bash
docker login docker.pkg.github.com -u <USERNAME> -p <GITHUB_TOKEN>
docker build -t  docker.pkg.github.com/kyberorg/honeypot/honeypot:latest .
# make container
docker push docker.pkg.github.com/kyberorg/honeypot/honeypot:latest
# make push
```

Pull and Run:

```bash
docker pull docker.pkg.github.com/kyberorg/honeypot/honeypot:latest
docker run docker.pkg.github.com/kyberorg/honeypot/honeypot:latest
```


### Setup new SonarCloud Project

- On _SonarCloud_:
    - Click _Plus_ Sign in Upper Right Corner
    - _Analyze New Project_
    - Click _GitHub app configuration_ link
    - Configure SonarCloud
    - Select Repository and Save
    - Go Back to Analyze Project
    - Tick Newly Added Repository
    - Click Set Up
    - Click Configure with Travis
    - Copy the Command to Encrypt the Travis Token
    - Run `travis encrypt --com <TOKEN_YOU_COPPIED>`
    - Populate the `secure` Field in `.travis.yml` with outputted string
    - Follow steps to populate your `sonar-project.properties`
    - Push
- On Travis CI:
    - Set `DOCKER_USERNAME`
    - Set `DOCKER_PASSWORD` to Your GitHub Registry Token

### Setup CodeClimate
- Go to <https://codeclimate.com/github/repos/new>
- Add Repository
- Go to Test Coverage Tab
- Copy Test Reporter ID
- Go to Travis and Open Settings for Your Repository
- Add Environment Variable: name: `CC_TEST_REPORTER_ID`, value: _Copied from CodeClimate_
