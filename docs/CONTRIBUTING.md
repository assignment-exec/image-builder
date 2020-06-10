# Contributing to Image Builder

The following is a set of guidelines for contributing to Image Builder and its packages, which are a part of [Assignment-Exec Organization](https://github.com/assignment-exec) on GitHub.

## Submit Changes or New Features

### Clone the source code
- To make changes to [image-builder](https://github.com/assignment-exec/image-builder) you need to first fork the repository and have the forked repository checked out locally.
- After forking, copy the https url to clone the repository and use it in the command given below.
```commandline
$ git clone <https_url_of_repo>
$ cd image-builder
```

### Make changes in a new branch
Each change must be made in a separate branch, created from `develop` branch. Use git commands to create a branch and add changes.
```commandline
$ git checkout -b <branch_name>
$ [edit files...]
$ git add [files...]
```

### Commit changes
- Write a clear log message for your commits. For small changes simple one-line messages are fine, but for big changes a detailed description should be provided.
```commandline
git commit
```
- First line in the commit messages should be a brief summary of the commit.
- Follow the first line by a blank line and then the main description about the changes.
- An ideal case would be that the commits are atomic i.e one feature or change is added per commit.

### Add Unit tests
- For every change unit test cases should be added if possible. All the tests should be passed before submitting the change for review.
- See [details](https://golang.org/pkg/testing/) for writing unit tests in golang.

### Send changes for review
- Once the change is ready send a pull request to the `develop` branch of [image-builder](https://github.com/assignment-exec/image-builder/pull/new/develop) with a list of things you've done.
- See [pull requests](http://help.github.com/pull-requests/) to know more about it. 


Thank you for taking time to contribute !