# wikiapp
Web application example in Go

As part of the Linux Software Foundation trainings, the LFD254 course provides the rsvpapp example to demonstrate different technologies around containers (Docker, docker-compose, etc.). Instead of simply re-using the provided "rsvpapp" application, I built this simple wiki application in Go and created the necessary Dockerfile and docker-compose.yml files.

This wiki app is adapted from https://golang.org/doc/articles/wiki/. Plus, it uses mongodb to store the content of the pages instead of files.

Update the README to trigger the Jenkins Pipeline
