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
