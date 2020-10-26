# Release process

## Preparation

Gather closed issues since last release and prepare the announcement
text.

## Releasing
We will assume we are releasing version `x.y.z`

```shell script
git pull 
git checkout master
git tag vx.y.z
git push origin vx.y.z
```

Automation should create a release for you. Edit the release text with
the announcement.


## Notes about docker images

Docker images are rolling release