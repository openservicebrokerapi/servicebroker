# Platform Compatibility for OSBAPI


| Release or Feature | Introduced | Deprecated Cloud Foundry | Kubernetes |
| --- | -- | --- | --- |
| Authentication - Basic | pre-v2.10 | | ✔️ | ✔️ |
| Authentication - Opaque Tokens | pre-v2.10 | v2.14 | - | ✔️ |
| Binding - Credentials | v2.10 | | ✔️ | ✔️ |
| Binding - Log Drain | v2.10 | | ✔️ | - |
| Binding - Route Services | v2.10 | | ✔️ | - |
| Binding - Volume Services | v2.10 | | ✔️ | - |
| Plan Bindable | v2.11 | | ✔️ | ✔️ |
| [*v2.11*](release-notes.md#v211) | Aug 11, 2016 | | ✔️ | ✔️ |
| `context` for PUT,PATCH instance | v2.12 | | ✔️ | ✔️ |
| [*v2.12*](release-notes.md#v212) | June 13, 2017 | | ✔️ | ✔️ |
| `schemas` in catalog | v2.13 | | ✔️ | ✔️ |
| `context` for PUT binding | v2.13 | | ✔️ | ✔️ |
| `volume_mounts` bindings | v2.13 | | ✔️ | - |
| `originating identity` header | v2.13 | | ✔️ | ✔️ |
| [*v2.13*](release-notes.md#v213) | Sep 27, 2017 | | ✔️ | ✔️ |
| k8s `context.clusterid` | v2.14 | | n/a | ✔️ |
| allow periods and uppercase letters in name fields | v2.14 | | - | ✔️ |
| Opaque Bearer Token Authentication | v2.14 | | - | ✔️ |
| Async Bindings | v2.14 | | - | ✔️ |
| GET endpoint for Service Instance | v2.14 | | ✔️ | - |
| GET endpoint for Service Binding | v2.14 | | ✔️ | ✔️ |
| Authentication - Bearer Tokens | v2.14 | | - | ✔️ |
| *v2.14* | _TBD_ | | - | - |
| `dashboard_url` is updatable on Service Instance Update | [_OPEN_](https://github.com/openservicebrokerapi/servicebroker/pull/437) | | - | - |
| Generic Actions | [_OPEN_](https://github.com/openservicebrokerapi/servicebroker/pull/431) | | - | - |
| Generic Broker Actions | _OPEN_ | - | - |
| JSON Schema Endpoint | [_OPEN_](https://github.com/openservicebrokerapi/servicebroker/pull/402) | | - | - |
| JSON Schema for Responses | [_OPEN_](https://github.com/openservicebrokerapi/servicebroker/pull/392) | | - | - |
| Deprecated Classes and Plans | [_OPEN_](https://github.com/openservicebrokerapi/servicebroker/pull/504) | | - | - |
