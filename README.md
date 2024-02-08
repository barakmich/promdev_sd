# promdev_sd
Prometheus Service Discovery for local/dev environments using HTTP_SD

## Basic Idea

A poor man's [Prometheus](https://prometheus.io) service discovery.

I'm running [`tilt`](https://tilt.dev) or `docker-compose` or similar as I'm developing my code, and want to hook into an org-wide Prometheus monitoring solution. These are very ephemeral addresses to scrape, but there's a network route from the scraper to my development box. I don't want to add my development box directly to the scrape list. I need a service discovery mechanism that's lightweight and temporary.

## Plan

- Organization runs a `promdev_server` which exposes:
  - `/register` -- for new scrape targets to heartbeat against
  - `/discovery` -- which is a HTTP_SD target added to the Prometheus config (either statically or because the `promdev_server` is running somewhere)
- Developer has options:
  - Runs `promdev_reporter` as a process (first goal)
    - It comes up and down with my development servers and simply heartbeats the static set of ports and `/metrics` endpoints up the chain
  - Integrates `promdev_sd` as a library (todo!)
    - The process itself, in development mode, heartbeats its own `/metrics` endpoint against the target `promdev_server`

## Goal
- Prometheus discovers the temporary development servers and scrapes them (with all appropriate tags).
- Developer sees metrics show up in the organization's dashboards (Grafana, et al) tagged for themselves.
- `edit->build->run some local calls->see metrics->edit again`, shut down for the day, start up again, your dashboard is there for you
- Going to staging/prod? Then you'll have the same metrics endpoints scraped by the orchestrator via its service discovery (eg, Kubernetes)
