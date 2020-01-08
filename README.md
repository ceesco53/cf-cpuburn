cf-cpuburn lets you use 100% of all available cores, useful when stress-testing.
Download cpuburn from http://patrickmylund.com/projects/cpuburn/
See http://patrickmylund.com/projects/cpuburn/ for more information.

== About cf-cpuburn (feel the burn edition)

This app converts cpuburn into a Cloud Foundry native app.

**NOTE before you push**
Be sure to check https://github.com/cloudfoundry/go-buildpack/releases and check the GOVERSION in the manifest matches with a version packaged with go_buildpack.

== Installation

## Run in the Cloud

git clone https://github.com/ceesco53/cf-cpuburn.git
cd cf-cpuburn
cf push

More instructions to follow

## things to work on:
Add an endpoint to stop/start the burn
