## TODO
- [x] readme
- [x] logic for packs
  - [x] reivisit pack logic eventually
- [x] tests for pack logic
- [x] local api setup
- [x] setup terraform for infra
- [x] CI/CD
- [x] lambda api entry point setup

STRETCH
- [x] ability to manage pack sizes
- [ ] 'thoughts' doc outlining my thinking
- [ ] frontend
- [ ] deploy frontend with infra

-----------------

TODO
- [ ] GetShippingPacksForOrder
  - [ ] initial 'greedy' algo attempt
  - [ ] rework in DP way
  - [ ] testing
- [ ] Repo structure
  - [ ] backend
    - [ ] initially setup for local & lambda seperate
    - [ ] docker pattern for aws different
  - [ ] terraform
    - [ ] bootstrap
    - [ ] ci/cd setup
- [ ] Improvements for future
  - [ ] auth for api
  - [ ] frontend

-----------------

# Thoughts

I wanted to include some details related to my thought process & findings I came across whilst working on this submission. The contents are as follows:

- [Thoughts](#thoughts)
  - [`GetShippingPacksForOrder`](#getshippingpacksfororder)
    - [Attempt 1](#attempt-1)
    - [Attempt 2](#attempt-2)
  - [Repo Structure](#repo-structure)
    - [`backend`](#backend)
      - [testing](#testing)
    - [`infra`](#infra)
  - [Improvements to be made](#improvements-to-be-made)
    - [API Auth](#api-auth)
    - [Frontend](#frontend)

## `GetShippingPacksForOrder`

### Attempt 1

[Link to initial solution](https://github.com/saddiqs1/gymshark-saddiqs1/pull/1/changes#diff-79df022a1973f733218d6dc0d856909acbdf227eec9d0ac9c504610380df56b3R4)

Initially, I implemented the first solution that popped into my head. My thought process was as follows:
- Sort the `packSizes` in descending order
- Loop through each `packSize`
- If the items left to package were larger than the `packSize`, use as many packs as you can to cover off the items
- Repeatedly do this until there are no more items to cover

This seemed like an obvious solution, and worked for the most part. However I ran into some edge cases pretty quickly (with the help of the [tests](https://github.com/saddiqs1/gymshark-saddiqs1/pull/1/changes#diff-9247afe83db5e0626cabd3265179a17392b74a0cdf0ddb4ea33fff4f3e8e2ca1R42) that I setup). I tried to cover off some basic edge cases with conditional statements in the `packSizes` loop, but it was not a clean solution, and it was hardcoded around the initial list of `packSizes` provided in the exercise. 

I knew this would need revisiting, but I wanted to focus on setting up everything else in the repo first (see [Repo Structure](#repo-structure)) before coming back to it at a later time.

### Attempt 2

This time round, I had already implemented the configurable `packSizes`, so I needed to keep that in mind. I also knew that the initial logic was flawed. Take the following scenario for example:
- `itemsOrdered` = 8
- `packSizes` = [3, 4, 6]
- The expected result would be: `4x2`, however the initial logic I had before would return `6x1, 3x1`. This violates the rule 2 of the exercise (i.e. sending out too many items)

With this in mind, I approached the problem from a different point of view, and came up with the following:
- It is possible to work out the 'worse case scenario', by fulfilling the order just using the smallest `packSize` available. This can be used as a maximum amount of items (i.e. `maxShippingPacksTotalItems`) we should send out
- Inversely, we know that the 'best possible scenario' is that we send out the amount of items ordered (i.e. `itemsOrdered`). This can be used as a minimum amount of items we aim to send out
- We can loop through every number from 1 --> `maxShippingPacksTotalItems`, and calculate the best combination of packs to ship for each item count
  - Within this step, by looping from 1 --> `maxShippingPacksTotalItems`, you can leverage pre-calculated combinations to make the calculations easier
- Once we've calculated every possible combination of packs to be sent out, loop through the minimum amount of items (`itemsOrdered`) to the maximum amount of items (`maxShippingPacksTotalItems`), and the first combination we find will be the least amount of items we are able to send with a valid combination of `packSizes`.

e.g. [link to diagram](https://excalidraw.com/#json=mC5zwll_ctEn9xBmvdCts,R4wuHWDhLcr6YhfNYOM-Kg)
![shippingpacks logic](shippingpacks-logic.png)


## Repo Structure

### `backend`

#### testing

### `infra`

## Improvements to be made

### API Auth

### Frontend
