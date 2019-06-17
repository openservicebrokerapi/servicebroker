# Platform Compatibility for OSBAPI

| Release or Feature | Introduced | Cloud Foundry | Kubernetes |
| --- | --- | --- | --- |
| `credentials` binding | v2.10 | ✔️ | ✔️ |
| `syslog_drain` binding | v2.10 | ✔️ | - |
| `route_forwarding` binding | v2.10 | ✔️ | - |
| `volume_mounts` binding | v2.10 | ✔️ | - |
| [*v2.10*](release-notes.md#v210) | Early 2016 | ✔️ | ✔️ |
| Bindable and non-bindable plans | v2.11 | ✔️ | ✔️ |
| [*v2.11*](release-notes.md#v211) | Nov 15, 2016 | ✔️ | ✔️ |
| `context` for creating and updating a Service Instance | v2.12 | ✔️ | ✔️ |
| [*v2.12*](release-notes.md#v212) | June 13, 2017 | ✔️ | ✔️ |
| `schemas` in catalog | v2.13 | ✔️ | ✔️ |
| `context` for creating a Service Binding | v2.13 | ✔️ | ✔️ |
| `originating identity` header | v2.13 | ✔️ | ✔️ |
| Opaque Bearer Token Authentication | v2.13 | - | ✔️ |
| [*v2.13*](release-notes.md#v213) | Sep 27, 2017 | ✔️ | ✔️ |
| GET endpoint for Service Instances | v2.14 | ✔️ | - |
| GET endpoint for Service Bindings | v2.14 | ✔️ | ✔️ |
| Async Bindings | v2.14 | ✔️ | ✔️ |
| [*v2.14*](release-notes.md#v214) | July 24, 2018 | ✔️ | ✔️ |
| `endpoints` information in Service Bindings used to setup network connectivity | v2.15 | - | - |
| `maintenance_info` used for Service Instance upgrades | v2.15 | ✔️ | - |
| `plan_upgradeable` field in Service Plan object is used to determine if a plan change can be performed | v2.15 | ✔️ | - |
| `maximum_polling_duration` field in Service Plan object is adhered to | v2.15 | ✔️ | - |
| `Retry-After` response header from calls to last operation is adhered to | v2.15 | ✔️ | - |
[ [*v2.15*](release-notes.md#v215) | June 25, 2019 | ✔️ | ✔️ |
