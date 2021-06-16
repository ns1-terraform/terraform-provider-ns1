## 1.10.3 (June 15, 2021)
ENHANCEMENTS:
* Add more verbose logging output for failed requests [#160](https://github.com/ns1-terraform/terraform-provider-ns1/pull/160)
* Update to documentation to reflect proper usage for monitoring datafeeds [#154](https://github.com/ns1-terraform/terraform-provider-ns1/pull/154) 
BUG FIXES:
* Correctly coerce `test_id` config value for datafeed resource used with `thousandeyes` `datasource` [#163](https://github.com/ns1-terraform/terraform-provider-ns1/pull/163)
* Change `notify_failback` field in `monitoringjob` resource default to `true` to match default in api [#161](https://github.com/ns1-terraform/terraform-provider-ns1/pull/161)

## 1.10.2 (May 21, 2021)
ENHANCEMENTS:
* Updates ns1-go dependency to add handling of rate limitting when API returns 4xx error [#159](https://github.com/ns1-terraform/terraform-provider-ns1/pull/159).

## 1.10.1 (April 27, 2021)
BUG FIXES:
* Resolves issue with missing value for Key attribute when creating an apikey joined to a team [#158](https://github.com/ns1-terraform/terraform-provider-ns1/pull/158).

## 1.10.0 (April 22, 2021)
ENHANCEMENTS:
* Adds DS record support [#157](https://github.com/ns1-terraform/terraform-provider-ns1/pull/157).

## 1.9.4 (March 23, 2021)
NOTES:
* Updates docs to clarify `key` is an Attribute and not an Argument [#150](https://github.com/ns1-terraform/terraform-provider-ns1/pull/150).

## 1.9.3 (March 4, 2021)
BUG FIXES:
* Adds missing `account_manage_ip_whitelist` permission [#148](https://github.com/ns1-terraform/terraform-provider-ns1/pull/148).

## 1.9.2 (February 26, 2021)
BUG FIXES:
* Values for `tls_skip_verify` are coerced correctly [#146](https://github.com/ns1-terraform/terraform-provider-ns1/pull/146). Thanks @zahiar!

## 1.9.1 (January 7, 2021)
BUG FIXES:
* Values for IPv6 monitoring job configs are coerced correctly [#141](https://github.com/ns1-terraform/terraform-provider-ns1/pull/141).

## 1.9.0 (September 8, 2020)
FEATURES:
* **New Data Source:** `ns1_record` [#137](https://github.com/ns1-terraform/terraform-provider-ns1/pull/137). Thanks to @zahiar!

## 1.8.6 (August 31, 2020)
ENHANCEMENTS:
* Add additional config field to monitoring job configuration

## 1.8.5 (August 13, 2020)
BUG FIXES:
* Resolves issue with config maps returning floats sometimes

## 1.8.4 (June 24, 2020)
BUG FIXES:
* Resolves an issue where changes involving feed pointers in record answer metadata were not detected ([124](https://github.com/terraform-providers/terraform-provider-ns1/pull/124))

## 1.8.3 (May 21, 2020)
BUG FIXES:
* Resolves issues on record filter and meta fields around boolean values not properly being converted to strings ([123](https://github.com/terraform-providers/terraform-provider-ns1/pull/123)).

## 1.8.2 (May 01, 2020)
NOTES:
* Clarify rate limit documentation ([121](https://github.com/terraform-providers/terraform-provider-ns1/pull/121))
* Replace examples in README with blurb pointing to docs ([120](https://github.com/terraform-providers/terraform-provider-ns1/pull/120))

## 1.8.1 (April 09, 2020)
ENHANCEMENTS
 * Change username validation regex to match validation used by NS1 API.([#119](https://github.com/terraform-providers/terraform-provider-ns1/pull/119))

## 1.8.0 (March 19, 2020)

ENHANCEMENTS:
* support for pulsar metadata in record answers ([#116](https://github.com/terraform-providers/terraform-provider-ns1/issues/116))

## 1.7.1 (February 25, 2020)

BUG FIXES:
* Bump ns1-go SDK version to v2.2.1 - resolves an issue with ASNs causing
  panics ([#113](https://github.com/terraform-providers/terraform-provider-ns1/pull/113)).
* Fix for IP Prefix ordering - don't show a change when order differs ([#112](https://github.com/terraform-providers/terraform-provider-ns1/pull/112)).

ENHANCEMENTS:
* Validate username field in the provider, so issues with usernames are caught
  in the "plan" stage ([#115](https://github.com/terraform-providers/terraform-provider-ns1/pull/115)).

## 1.7.0 (January 28, 2020)

NOTES:
* The `short_answers` attribute on `ns1_record` has had a deprecation warning added to it and will be deprecated in a future release ([#102](https://github.com/terraform-providers/terraform-provider-ns1/pull/102)).
* The project has been tagged as under "active development", in accordance with NS1 standards around public facing repositories ([#109](https://github.com/terraform-providers/terraform-provider-ns1/pull/109)).

ENHANCEMENTS:
* Support for DDI permissions on teams, users, and API keys has been added,
and can be enabled via the new `enable_ddi` configuration option on the provider ([#105](https://github.com/terraform-providers/terraform-provider-ns1/pull/105)).
* Added IP Whitelist support for teams, users, and AIP keys ([#105](https://github.com/terraform-providers/terraform-provider-ns1/pull/105)).
* Clarified documentation for IPv4 only fields ([#108](https://github.com/terraform-providers/terraform-provider-ns1/pull/108)).

## 1.6.4 (January 06, 2020)

IMPROVEMENTS:
* Updated permissions behavior on user and API key resources to accurately show `terraform plan` differences when the user or key is part of a team and updated documentation accordingly ([#100](https://github.com/terraform-providers/terraform-provider-ns1/pull/100))
* Switched to the Terraform standalone SDK ([#101](https://github.com/terraform-providers/terraform-provider-ns1/pull/101))
* Update resource state management to properly handle disappearing resources ([#99](https://github.com/terraform-providers/terraform-provider-ns1/pull/99))

## 1.6.3 (December 16, 2019)

IMPROVEMENTS:
* Add validation to the zone and domain fields on a record to more clearly indicate invalid inputs containing leading or trailing dots ([#97](https://github.com/terraform-providers/terraform-provider-ns1/pull/97))

## 1.6.2 (December 05, 2019)

ENHANCEMENTS:
* Support URLFWD records ([#96](https://github.com/terraform-providers/terraform-provider-ns1/issues/96))
* Add a "clean" rule to Makefile ([#89](https://github.com/terraform-providers/terraform-provider-ns1/issues/89))

## 1.6.1 (November 13, 2019)

BUG FIXES:
* fix interaction with the `autogenerate_ns_record` flag that was making terraform think a clean resource was dirty ([#85](https://github.com/terraform-providers/terraform-provider-ns1/issues/85))

ENHANCEMENTS:
* docs and example for using `autogenerate_ns_record`. ([#83](https://github.com/terraform-providers/terraform-provider-ns1/issues/83))
* minor improvements to some error messages in tests.
* improve docs around the ordering requirement for zone regions.
* improve docs around provider arguments and environment variables.

IMPROVEMENTS:
* Add a configuration option to the provider to use an alternate strategy to avoid rate limit errors. ([#88](https://github.com/terraform-providers/terraform-provider-ns1/issues/88))

## 1.6.0 (October 16, 2019)

BUG FIXES:
* Pick up a divide by zero fix in the SDK, when rate limiting. Should wait around when hitting limits rather
  than falling over, and shouldn't require limiting parallelism to avoid 429 errors. ([#74](https://github.com/terraform-providers/terraform-provider-ns1/issues/74))
 * We were explicitly sending defaults for `port` and `notify` fields in `secondaries`, now they are implicit.
   Sending the default `port` prevented using IP ranges. ([#82](https://github.com/terraform-providers/terraform-provider-ns1/issues/82))
 
ENHANCEMENTS:
* Allow secondary zone -> primary zone in-place. Old behavior was to force a new resource (DELETE/PUT) on any
  change to secondary. Now the only one that requires a new resource is when a zone _becomes_ a secondary. See
  the note in docs. ([#75](https://github.com/terraform-providers/terraform-provider-ns1/issues/75))
* Support for CAA records. ([#78](https://github.com/terraform-providers/terraform-provider-ns1/issues/78))
  
DEPRECATION:
* We've bumped the CI tests to run under Go 1.12. Provider still works with 1.11, but we're developing on and
  targeting 1.12+ ([#76](https://github.com/terraform-providers/terraform-provider-ns1/issues/76))

## 1.5.3 (October 02, 2019)

BUG FIXES:
* Disallow setting SOA fields (refresh, retry, expiry, and nx_ttl) on secondary zones. An API bug allowed
  these fields to be set on "create", the API now discards any settings to these fields and sets them to default
  values. These fields are now marked as "ConflictsWith" for secondary zones. If you were doing this and tf complains
  or the plan becomes dirty, the solution is to ensure the values are correctly set to the defaults, and let these
  fields be computed. ([#71](https://github.com/terraform-providers/terraform-provider-ns1/issues/71))

ENHANCEMENTS:

* Support toggling DNSSEC on zones: requires account to have DNSSEC permission (this is managed by support) ([#70](https://github.com/terraform-providers/terraform-provider-ns1/issues/70))
* Zone DNSSEC data_source: this data source has the DNSSEC info for a zone with DNSSEC enabled ([#72](https://github.com/terraform-providers/terraform-provider-ns1/issues/72))

## 1.5.2 (September 20, 2019)

ENHANCEMENTS:

* Support outgoing transfer a.k.a. primary zone ([#65](https://github.com/terraform-providers/terraform-provider-ns1/issues/65))
* Add option to enable request body logging via env var ([#67](https://github.com/terraform-providers/terraform-provider-ns1/issues/67))

IMPROVEMENTS:

* acc tests: Randomize zone names to help prevent collisions ([#64](https://github.com/terraform-providers/terraform-provider-ns1/issues/64))
* Ignore order of location fields (comma sep strings) in record regions block [[#68](https://github.com/terraform-providers/terraform-provider-ns1/issues/68)] 
* Correct and improve docs around "regions" in record resource ([#69](https://github.com/terraform-providers/terraform-provider-ns1/issues/69))

## 1.5.1 (August 30, 2019)

IMPROVEMENTS:

* General Documentation updates: Flesh out examples and attributes across the board. ([#62](https://github.com/terraform-providers/terraform-provider-ns1/issues/62))
* Add known issues/roadmap section to main README ([#62](https://github.com/terraform-providers/terraform-provider-ns1/issues/62))
* resource/ns1_record: Document import and `meta` arguments. ([#62](https://github.com/terraform-providers/terraform-provider-ns1/issues/62))

## 1.5.0 (July 31, 2019)

ENHANCEMENTS:

* resource/notifylist: Add support for all notifier types currently supported by SDK. [#59](https://github.com/terraform-providers/terraform-provider-ns1/pull/59)
* resource/ns1_zone: Add `additional_primaries` argument. Add documentation for all arguments and attributes. ([#60](https://github.com/terraform-providers/terraform-provider-ns1/issues/60))
* datasource/ns1_zone: Add `additional_primaries` attribute. Add documentation for all arguments and attributes. ([#60](https://github.com/terraform-providers/terraform-provider-ns1/issues/60))

IMPROVEMENTS:
* Updates ns1-go dependency to latest version ([#60](https://github.com/terraform-providers/terraform-provider-ns1/issues/60))


## 1.4.1 (July 04, 2019)

IMPROVEMENTS:
* Update ns1-go dependency to latest version to fix bug when passing multiple `ip_prefixes` as comma delimited string [#57](https://github.com/terraform-providers/terraform-provider-ns1/pull/57)
* Update Terraform dependency to v0.12.3 [#58](https://github.com/terraform-providers/terraform-provider-ns1/pull/58)

## 1.4.0 (May 13, 2019)

IMPROVEMENTS:

* Update Terraform dependency to v0.12.0-rc1 [#55](https://github.com/terraform-providers/terraform-provider-ns1/pull/55)

## 1.3.1 (April 30, 2019)

IMPROVEMENTS:

* Update ns1-go dependency to latest version [#54](https://github.com/terraform-providers/terraform-provider-ns1/pull/54). Thanks to @glennslaven!

## 1.3.0 (April 09, 2019)

BUG FIXES:

* resource/user: Force new user on modification of username [#28](https://github.com/terraform-providers/terraform-provider-ns1/pull/28).  Thanks to @jamesgoodhouse!
* resource/record: Sort regions inside records to ensure deterministic comparison between configuration and current state [#49](https://github.com/terraform-providers/terraform-provider-ns1/pull/49). Regions in a record's region list will now need to be sorted alphanumerically by name, otherwise a modification will be detected when none actually exists. Thanks to @bparli!

## 1.2.0 (March 26, 2019)

FEATURES:

* **New Data Source:** `ns1_zone` [#47](https://github.com/terraform-providers/terraform-provider-ns1/pull/47)

IMPROVEMENTS:

* Migrate to Go Modules [#48](https://github.com/terraform-providers/terraform-provider-ns1/pull/48))
* Refactor acceptance test fixtures to Terraform 0.12 syntax [#50](https://github.com/terraform-providers/terraform-provider-ns1/pull/50)
* Update website and examples to Terraform 0.12 syntax [#51](https://github.com/terraform-providers/terraform-provider-ns1/pull/51)
* Update ns1-go module latest version [#51](https://github.com/terraform-providers/terraform-provider-ns1/pull/51)

## 1.1.2 (February 06, 2019)

BUG FIXES:

* resource/record: Don't try to convert config values to ints

## 1.1.1 (February 01, 2019)

BUG FIXES:

* resource/zone: Make `networks` field computed

## 1.1.0 (January 08, 2019)

* Add support for short_answers in record resources

## 1.0.0 (January 25, 2018)

* Metadata support implemented for records, answers, and regions
* Small bugfixes

## 0.1.0 (June 21, 2017)

* NS1 Support for Terraform implemented

NOTES:

* Same functionality as that of Terraform 0.9.8. Repacked as part of [Provider Splitout](https://www.hashicorp.com/blog/upcoming-provider-changes-in-terraform-0-10/)
