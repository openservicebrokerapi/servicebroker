# Contributing

All contributors sending Pull Requests (PRs) must have a Contributor
License Agreement on file as either an
[individual](https://www.cloudfoundry.org/pdfs/CFF_Individual_CLA.pdf)
or via their
[employer](https://www.cloudfoundry.org/pdfs/CFF_Corporate_CLA.pdf).

## Minor Changes

Minor change proposals to the specification, changes such as editorial bugs
or enhancements that do not modify the semantics of the specification or
syntax of the API, can be suggested via a Github Issue or Pull Request (PR).

If there is a need for some discussion around how best to
address the concern then opening an Issue, prior to doing the work to develop
a PR, would be best.  These minor issues do not need a formal review, per the
[process](#prissue-review-process) described below, rather the issue
should be used to come to a conscensus around a PR that will eventually be
submitted.

If the proposed change is non-controversial (e.g. a typo) then a PR can be
submitted directly without opening an issue.

Either way, once a PR is submitted it will be reviewed per the
[process](#prissue-review-process) described below.

## Design Change Proposals

New design proposals to the API spec should be submitted by opening a
Github issue with a link to a Google Doc containing the proposal. Proposals
should focus primarily on motivation, problem statement, and use cases before
suggesting a solution or change.  Collaboration on the design, fleshing out of
use cases, etc can occur as comment discussions in the Google Doc, as well
as on our weekly calls.

For changes largely impacting the SB API actors (platforms, service authors,
etc...), it is recommended to solicit feedback from these actors and leave
enough time (say 2 weeks) for feedback to be provided, and for the
potentially received objections/suggestions to be handled.

Once the design has been finalized in the Google Doc, the proposed set of
changes to the specification should be made available for review. This could
be done by pointing to a branch in a Github repo with the proposed edits.
Reviewers, or implementers of the feature, can then easily see the exact
changes being proposed and comment on them within that branch.

Before a formal PR for the changes is created, the feature needs to be
validated through a full implementation by at least one platform that is
easily accessible to reviewers. Note that in many cases this would include
support for the proposed changes by one or more service brokers
(potentially samples). How each platform makes this new feature available
will vary between platforms. It is strongly recommended that as part of the
design review process platforms that plan on implementing the feature
get confirmation from the reviewers that their planned mechanism for providing
access to the implementation of the feature is acceptable.

It is expected that during this implementation phase there will be changes
made to the design and proposed specification edits to accurately represent
the current status of the proposal.

Once support for the change has been implemented, a PR should be sent to
this repo for the change to the broker API specification. By this point,
the API interactions should be well understood and there should be no
technical surprises; we expect the only discussion necessary on PRs to be
for wordsmithing and formatting.

## PR/Issue Review Process

All proposals (either Pull Requests or Issues) will follow the process
described below:

- All proposals must be available for review for no less than one week before
  they are approved. This will provide each dedicated committer enough time
  to review the proposal without unnecessarily delaying forward progress.
  Any non-trivial edit to the proposal (e.g. edits larger than typos) will
  reset the clock.
- Any dedicated committer can veto (via a "NOT LGTM" comment in the proposal).
  The comment must include the reasoning behind the veto.
  It is then expected that the dedicated committers will discuss the concerns
  and determine the next step for proposal - either close/reject the proposal
  or address the concerns raised such that the "NOT LGTM" can be rescinded.
- A proposal requires at least 4 "LGTM" comments from at least 4 different
  organizations to be approved.
- Once a "design change" issue is approved, it will be tagged with an
  "proposal finalized" label.  This indicates that it is ready to be
  implemented by a platform developer, see the [process](#contributing) above.
- Once a Pull Request is approved, it will be merged into the 'master' branch.

## Spec Release Process

Any member of the OSBAPI PMC can request that the accepted spec changes on
master are released in a new version of the spec. They will do this by
creating a new PR with the tile of the proposed release. For example,
"Release Proposal: v2.20".

### Prepare a PR

- In a fork create a new branch called "v$major.$minor".
- Update [README](README.md)'s "Latest Release" section and link to branch URL.
- Update [RELEASE-NOTES](release-notes.md) with details of the changes included
  in the release.

### Review Process

- All release proposals must be available for review for no less than one
  week before they are approved. This will provide each dedicated committer
  enough time to review the release proposal without unnecessarily delaying
  forward progress.
- Any dedicated committer can veto (via a "NOT LGTM" comment in the proposal).
  The comment must include the reasoning behind the veto.
  It is then expected that the dedicated committers will discuss the concerns
  and determine the next step for release proposal
	- either close/reject the proposal or address the concerns raised such that
	the "NOT LGTM" can be rescinded.
- A release proposal requires at least 4 "LGTM" comments from at least
  4 different organizations to be approved.

### Once Approved

Once the release is approved, the following actions should be taken by
a PMC member:

- Merge PR into master.
- Create a branch from the SHA of the merge named "v$major.$minor".
- Notify mailing list of the new release!
