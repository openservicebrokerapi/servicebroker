# Contributing

New proposals to the API spec should be submitted by opening a Github issue with a link to a Google Doc containing the proposal. Proposals should contain motivation, use cases, and detailed example of proposed changes to API interactions. Collaboration on the design, fleshing out of use cases, etc can occur as comment discussions in the Google Doc, as well as on our weekly calls.

Once the design has been finalized in the Google Doc, support for the change must be implemented in the service controller of a platform easily accessible to users. Currently Cloud Foundry is the most appropriate platform for this, but we look forward to other platforms adding support for the broker API soon, and one of those platforms could take the initiative of implementing support for the feature earlier. 

Once support for the change has been implemented, a PR should be to this repo for the change to the broker API specificaion. By this point, the API interactions should be well understood and there should be no technical surprises; we expect the only discussion necessary on PRs to be for wordsmithing and formatting.

All contributors sending Pull Requests (PRs) must have a Contributor License Agreement 
on file as either an [individual](https://www.cloudfoundry.org/pdfs/CFF_Individual_CLA.pdf) 
or via their [employer](https://www.cloudfoundry.org/pdfs/CFF_Corporate_CLA.pdf).
