## [1.4.4](https://github.com/persianopencart/fleetssl-cpanel-new/compare/v1.4.3...v1.4.4) (2026-05-17)


### Reverts

* restore WHM API token to acl-1=all ([b1fbd17](https://github.com/persianopencart/fleetssl-cpanel-new/commit/b1fbd17b20618e281117b19e40e25ed1a6466d30))

## [1.4.3](https://github.com/persianopencart/fleetssl-cpanel-new/compare/v1.4.2...v1.4.3) (2026-05-17)


### Bug Fixes

* scope the WHM API token to least privilege instead of root ([975fa95](https://github.com/persianopencart/fleetssl-cpanel-new/commit/975fa95fac97fb8984280868998d635a2f6ea39c))

## [1.4.2](https://github.com/persianopencart/fleetssl-cpanel-new/compare/v1.4.1...v1.4.2) (2026-05-17)


### Bug Fixes

* create WHM API tokens with privileges and stop them accumulating ([aa92ba6](https://github.com/persianopencart/fleetssl-cpanel-new/commit/aa92ba6c27eacf75c434a55bb0d80d2a52c8baee))

## [1.4.1](https://github.com/persianopencart/fleetssl-cpanel-new/compare/v1.4.0...v1.4.1) (2026-05-17)


### Bug Fixes

* base wildcard DNS viability on the requested name ([4ad9c4c](https://github.com/persianopencart/fleetssl-cpanel-new/commit/4ad9c4cb62eb011d8725b11b98fba40527a37ae6))

# [1.4.0](https://github.com/persianopencart/fleetssl-cpanel-new/compare/v1.3.0...v1.4.0) (2026-05-16)


### Features

* smart wildcard detection and resilient multi-domain issuance ([b86de4b](https://github.com/persianopencart/fleetssl-cpanel-new/commit/b86de4b18b24a8c1ce45904678cd9496d2e9563f))

# [1.3.0](https://github.com/persianopencart/fleetssl-cpanel-new/compare/v1.2.2...v1.3.0) (2026-05-16)


### Bug Fixes

* auto-recover from a rejected WHM API token ([e867cab](https://github.com/persianopencart/fleetssl-cpanel-new/commit/e867cab7679f979b27484d2ee09c55105cc1361a))


### Features

* modernize the UI and drop the bundled front-end libraries ([1e3db75](https://github.com/persianopencart/fleetssl-cpanel-new/commit/1e3db759365cecfb4ce00b472d4f1baa3eceb9ee))

## [1.2.2](https://github.com/persianopencart/fleetssl-cpanel-new/compare/v1.2.1...v1.2.2) (2026-05-16)


### Bug Fixes

* modernize dependencies and toolchain, clearing 31 known CVEs ([5557927](https://github.com/persianopencart/fleetssl-cpanel-new/commit/55579276033f14041b3bc444eebd07984c870b2f))

## [1.2.1](https://github.com/persianopencart/fleetssl-cpanel-new/compare/v1.2.0...v1.2.1) (2026-05-16)


### Bug Fixes

* always fall back to HTTP-01 when DNS-01 validation fails ([9e91ecb](https://github.com/persianopencart/fleetssl-cpanel-new/commit/9e91ecb08a15cca6deffed2b31bae2aeb1acb5fe))

# [1.2.0](https://github.com/persianopencart/fleetssl-cpanel-new/compare/v1.1.0...v1.2.0) (2026-05-16)


### Features

* publish installable RPM and DEB packages with each release ([1e32b49](https://github.com/persianopencart/fleetssl-cpanel-new/commit/1e32b4982ef0f9858c6310d3bff650d23b1a80f5))

# [1.1.0](https://github.com/persianopencart/fleetssl-cpanel-new/compare/v1.0.0...v1.1.0) (2026-05-16)


### Features

* automatic DNS-01 -> HTTP-01 validation fallback ([fe1de7e](https://github.com/persianopencart/fleetssl-cpanel-new/commit/fe1de7ed6a1c01f4b788b96ef5f8940e7fbe5b5c))
* modern UI refresh with local-only assets ([61a05a7](https://github.com/persianopencart/fleetssl-cpanel-new/commit/61a05a772b27c1d26c764417d5230950854ff9a2))

# 1.0.0 (2026-05-16)


### Bug Fixes

* **ci:** correct semantic-release workflow indentation ([12bb601](https://github.com/persianopencart/fleetssl-cpanel-new/commit/12bb601866a44bf03ae334bd3bd25c7e7a11c175))


### Features

* initial public release setup ([1ba7846](https://github.com/persianopencart/fleetssl-cpanel-new/commit/1ba78468506abcfede1b82dbe48f4d15e0bd21d4))
