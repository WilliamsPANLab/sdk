# Flywheel SDK

An SDK for interaction with a remote Flywheel instance, in Golang, Python, and Matlab!

[![GoDoc](https://godoc.org/github.com/flywheel-io/sdk/api?status.svg)](https://godoc.org/github.com/flywheel-io/sdk/api)
[![Report Card](https://goreportcard.com/badge/github.com/flywheel-io/sdk)](https://goreportcard.com/report/github.com/flywheel-io/sdk)
[![Build status](https://circleci.com/gh/flywheel-io/sdk/tree/master.svg?style=shield)](https://circleci.com/gh/flywheel-io/sdk)

## Building

```bash
git clone https://github.com/flywheel-io/sdk workspace/src/flywheel.io/sdk
ln -s workspace/src/flywheel.io/sdk sdk

./sdk/make.sh
```

This builds the golang SDK.<br/>
For other languages, check out the [bridge readme](bridge).

## Testing

The simplest way to run the test suite is to install the [CircleCI runner](https://circleci.com/docs/2.0/local-jobs/#installation) and use it from the SDK folder:

```
circleci build
```

If you want to test manually, you can configure the test suite with these environment variables:

* `SdkTestKey`: Set this to an API key. Defaults to `localhost:8443:change-me`.
* `SdkTestMongo`: Set this to a mongo connection string. If not set, database tests are skipped.
* `SdkTestDebug`: Setting this to any value will cause each test to print an HTTP/1.1 representation of each request. Best used to debug a single failing test.

To run the integration test suite against a running API:

```bash
export SdkTestKey="localhost:8443:some-api-key"

./sdk/make.sh test
```

Or, to run a single test:

```bash
./sdk/make.sh test -run TestSuite/TestGetConfig
```

## Route Implementation Status

Route                                            | Golang  |  C++   | Python | Matlab
-------------------------------------------------|---------|--------|--------|--------
Get current user                                 | X       | X      | X      | X
&nbsp;                                           |         |        |        |
Get all users                                    | X       | X      | X      | X
Get user                                         | X       | X      | X      | X
Add user                                         | X       | X      | X      | X
Modify user                                      | X       | X      | X      | X
Delete user                                      | X       | X      | X      | X
&nbsp;                                           |         |        |        |
Get all containers of type                       | X       | X      | X      | X
Create container                                 | X       | X      | X      | X
Get container                                    | X       | X      | X      | X
Modify container                                 | X       | X      | X      | X
Delete container                                 | X       | X      | X      | X
Upload file to container                         | X       | X      | X      | X
Download file from container                     | X       | X      | X      | X
Add note to a container                          | X       | X      | X      | X
Upload tag to a container                        | X       | X      | X      | X
Get jobs that involve container                  |         |        |        |
&nbsp;                                           |         |        |        |
Set file attributes                              | X       | X      | X      | X
Set file info fields                             | X       | X      | X      | X
Replaces all file info fields                    | X       | X      | X      | X
Delete file info fields                          | X       | X      | X      | X
&nbsp;                                           |         |        |        |
Get all collections                              | X       | X      | X      | X
Get collection                                   | X       | X      | X      | X
Add collection                                   | X       | X      | X      | X
Modify collection                                | X       | X      | X      | X
Add session to collection                        | X       | X      | X      | X
Add acquisition to collection                    | X       | X      | X      | X
Get collection's sessions                        | X       | X      | X      | X
Get collection's acquisitions                    | X       | X      | X      | X
Get collection's session's acquisitions          | X       | X      | X      | X
Delete collection                                | X       | X      | X      | X
Add note to a collection                         | X       | X      | X      | X
&nbsp;                                           |         |        |        |
Create analysis                                  | X       |        |        |
Add note to analysis                             | X       | X      | X      | X
&nbsp;                                           |         |        |        |
Resolve path to route                            |         |        |        |
&nbsp;                                           |         |        |        |
Get all gears                                    | X       | X      | X      | X
Create gear                                      | X       | X      | X      | X
Get gear invocation                              | X       |        |        |
Suggest files for gear                           |         |        |        |
Delete gear                                      | X       | X      | X      | X
&nbsp;                                           |         |        |        |
Get a job                                        | X       | X      | X      | X
Get a job's logs                                 | X       | X      | X      | X
Append to a job's logs                           | X       |        |        |
Enqueue a job                                    | X       | X      | X      | X
Claim next pending job, and mark as running      | X       |        |        |
Modify job                                       | X       | X      | X      | X
&nbsp;                                           |         |        |        |
Get all batch jobs                               | X       | X      | X      | X
Get batch job                                    | X       | X      | X      | X
Propose batch job                                | X       |        |        |
Start batch job                                  | X       | X      | X      | X
Cancel batch job                                 | X       |        |        |
&nbsp;                                           |         |        |        |
Search                                           | X       | X      | X      | X
&nbsp;                                           |         |        |        |
Create bulk download ticket                      |         |        |        |
Get bulk download from tricket                   |         |        |        |
&nbsp;                                           |         |        |        |
Various upload strategies?                       |         |        |        |
Engine upload                                    |         |        |        |
&nbsp;                                           |         |        |        |
Declare a packfile upload to container           |         |        |        |
Upload to packfile                               |         |        |        |
Complete packfile and listen for progress        |         |        |        |
&nbsp;                                           |         |        |        |
Set project template                             |         |        |        |
Delete project template                          |         |        |        |
Recalculate project template compliance          |         |        |        |
Recalculate all project template compliance      |         |        |        |
&nbsp;                                           |         |        |        |
Get all devices                                  |         |        |        |
Get device                                       |         |        |        |
Get device statuses                              |         |        |        |
Get current device                               |         |        |        |
&nbsp;                                           |         |        |        |
Get site configuration                           | X       | X      | X      | X
Get site version                                 | X       | X      | X      | X
&nbsp;                                           |         |        |        |
_Delayed_                                        |         |        |        |
Regenerate current user's API key                |         |        |        |
Get user groups                                  |         |        |        |
Get containers for user?                         |         |        |        |
Clean out expired packfile progress              |         |        |        |
Scan for and fix disconnected jobs               |         |        |        |
Retry job                                        |         |        |        |
List groups with projects the user can access    |         |        |        |
Get schema                                       |         |        |        |
Get job stats (redesign on the horizon)          |         |        |        |
Get report                                       |         |        |        |
Get all gear rules                               |         |        |        |
Overwrite all gear rules                         |         |        |        |
&nbsp;                                           |         |        |        |
_Won't be implemented_                           |         |        |        |
List known sites (depreciated)                   |         |        |        |
Register a site (depreciated)                    |         |        |        |
Get current user avatar (no point)               |         |        |        |
Get user avatar (no point)                       |         |        |        |
Get all jobs (depreciated)                       |         |        |        |
Get job configuration (no point)                 |         |        |        |
