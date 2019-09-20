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
