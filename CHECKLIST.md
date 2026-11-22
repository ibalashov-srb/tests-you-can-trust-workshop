# A review checklist for tests you did not write

An agent gives you two hundred lines of tests. Ask these eight questions. If you cannot answer one, the test is not ready, even when it is green.

1. **Which rule in the spec does this test prove?** If there is no rule, the test only repeats the code.
2. **Was it seen red? Where is the proof?** Ask for the list of mutations: the line that was broken and the test that failed.
3. **Would it still pass if the behaviour was removed?** If yes, the test checks nothing.
4. **Does each "rejected" check have a control request?** The same request from a caller who is allowed must succeed. Without it, "rejected" is true for any error.
5. **Is the test data different from the defaults?** Code with a fixed constant passes a test that uses the default values.
6. **Does it check error values, not text or logs?** The text of a message changes; the error value is the contract.
7. **Are the mocks generated, and is the database real?** A mock written by hand stops matching its interface. A mocked database hides what only a real one can break.
8. **Does CI really run this test?** A test that does not run always looks green.
