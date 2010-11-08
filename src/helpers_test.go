// These tests verify the inner workings of the helper methods associated
// with gocheck.T.

package gocheck_test

import (
    "gocheck"
    "os"
)


var helpersS = gocheck.Suite(&HelpersS{})

type HelpersS struct{}

func (s *HelpersS) TestCountSuite(c *gocheck.C) {
    suitesRun += 1
}


// -----------------------------------------------------------------------
// Fake checker to verify the behavior of Assert() and Check().

type MyChecker struct {
    failCheck, noExpectedValue bool
    obtained, expected interface{}
}

func (checker *MyChecker) Name() string {
    return "MyChecker"
}

func (checker *MyChecker) ObtainedLabel() string {
    return "MyObtained"
}

func (checker *MyChecker) ExpectedLabel() string {
    return "MyExpected"
}

func (checker *MyChecker) NeedsExpectedValue() bool {
    return !checker.noExpectedValue
}

func (checker *MyChecker) Check(obtained, expected interface{}) bool {
    checker.obtained = obtained
    checker.expected = expected
    return !checker.failCheck
}


// -----------------------------------------------------------------------
// Tests for Check(), mostly the same as for Assert() following these.

func (s *HelpersS) TestCheckSucceedWithExpected(c *gocheck.C) {
    checker := &MyChecker{}
    testHelperSuccess(c, "Check(1, checker, 2)", true, func() interface{} {
        return c.Check(1, checker, 2)
    })
    if checker.obtained != 1 || checker.expected != 2 {
        c.Fatalf("Bad (obtained, expected) values for check (%d, %d)",
                 checker.obtained, checker.expected)
    }
}

func (s *HelpersS) TestCheckSucceedWithoutExpected(c *gocheck.C) {
    checker := &MyChecker{noExpectedValue: true}
    testHelperSuccess(c, "Check(1, checker)", true, func() interface{} {
        return c.Check(1, checker)
    })
    if checker.obtained != 1 || checker.expected != nil {
        c.Fatalf("Bad (obtained, expected) values for check (%d, %d)",
                 checker.obtained, checker.expected)
    }
}

func (s *HelpersS) TestCheckWithMissingExpected(c *gocheck.C) {
    checker := &MyChecker{failCheck: true}
    log := "helpers_test\\.go:[0-9]+ > helpers_test\\.go:[0-9]+:\n" +
           "\\.+ Check\\(obtained, MyChecker, >expected<\\):\n" +
           "\\.+ Missing expected value for the MyChecker checker\n\n"
    testHelperFailure(c, "Check(1, checker, !?)", nil, true, log,
                      func() interface{} {
        return c.Check(1, checker)
    })
}

func (s *HelpersS) TestCheckFailWithExpected(c *gocheck.C) {
    checker := &MyChecker{failCheck: true}
    log := "helpers_test\\.go:[0-9]+ > helpers_test\\.go:[0-9]+:\n" +
           "\\.+ Check\\(obtained, MyChecker, expected\\):\n" +
           "\\.+ MyObtained \\(int\\): 1\n" +
           "\\.+ MyExpected \\(int\\): 2\n\n"
    testHelperFailure(c, "Check(1, checker, 2)", false, false, log,
                      func() interface{} {
        return c.Check(1, checker, 2)
    })
}

func (s *HelpersS) TestCheckFailWithExpectedAndMessage(c *gocheck.C) {
    checker := &MyChecker{failCheck: true}
    log := "helpers_test\\.go:[0-9]+ > helpers_test\\.go:[0-9]+:\n" +
           "\\.+ Check\\(obtained, MyChecker, expected\\):\n" +
           "\\.+ MyObtained \\(int\\): 1\n" +
           "\\.+ MyExpected \\(int\\): 2\n" +
           "\\.+ Hello world!\n\n"
    testHelperFailure(c, "Check(1, checker, 2, msg)", false, false, log,
                      func() interface{} {
        return c.Check(1, checker, 2, "Hello", " world!")
    })
}

func (s *HelpersS) TestCheckFailWithoutExpected(c *gocheck.C) {
    checker := &MyChecker{failCheck: true, noExpectedValue: true}
    log := "helpers_test\\.go:[0-9]+ > helpers_test\\.go:[0-9]+:\n" +
           "\\.+ Check\\(obtained, MyChecker\\):\n" +
           "\\.+ MyObtained \\(int\\): 1\n\n"
    testHelperFailure(c, "Check(1, checker)", false, false, log,
                      func() interface{} {
        return c.Check(1, checker)
    })
}

func (s *HelpersS) TestCheckFailWithoutExpectedAndMessage(c *gocheck.C) {
    checker := &MyChecker{failCheck: true, noExpectedValue: true}
    log := "helpers_test\\.go:[0-9]+ > helpers_test\\.go:[0-9]+:\n" +
           "\\.+ Check\\(obtained, MyChecker\\):\n" +
           "\\.+ MyObtained \\(int\\): 1\n" +
           "\\.+ Hello world!\n\n"
    testHelperFailure(c, "Check(1, checker, msg)", false, false, log,
                      func() interface{} {
        return c.Check(1, checker, "Hello", " world!")
    })
}

// -----------------------------------------------------------------------
// Tests for Assert(), mostly the same as for Check() above.

func (s *HelpersS) TestAssertSucceedWithExpected(c *gocheck.C) {
    checker := &MyChecker{}
    testHelperSuccess(c, "Assert(1, checker, 2)", nil, func() interface{} {
        c.Assert(1, checker, 2)
        return nil
    })
    if checker.obtained != 1 || checker.expected != 2 {
        c.Fatalf("Bad (obtained, expected) values for check (%d, %d)",
                 checker.obtained, checker.expected)
    }
}

func (s *HelpersS) TestAssertSucceedWithoutExpected(c *gocheck.C) {
    checker := &MyChecker{noExpectedValue: true}
    testHelperSuccess(c, "Assert(1, checker)", nil, func() interface{} {
        c.Assert(1, checker)
        return nil
    })
    if checker.obtained != 1 || checker.expected != nil {
        c.Fatalf("Bad (obtained, expected) values for check (%d, %d)",
                 checker.obtained, checker.expected)
    }
}

func (s *HelpersS) TestAssertWithMissingExpected(c *gocheck.C) {
    checker := &MyChecker{failCheck: true}
    log := "helpers_test\\.go:[0-9]+ > helpers_test\\.go:[0-9]+:\n" +
           "\\.+ Assert\\(obtained, MyChecker, >expected<\\):\n" +
           "\\.+ Missing expected value for the MyChecker checker\n\n"
    testHelperFailure(c, "Assert(1, checker, !?)", nil, true, log,
                      func() interface{} {
        c.Assert(1, checker)
        return nil
    })
}

func (s *HelpersS) TestAssertFailWithExpected(c *gocheck.C) {
    checker := &MyChecker{failCheck: true}
    log := "helpers_test\\.go:[0-9]+ > helpers_test\\.go:[0-9]+:\n" +
           "\\.+ Assert\\(obtained, MyChecker, expected\\):\n" +
           "\\.+ MyObtained \\(int\\): 1\n" +
           "\\.+ MyExpected \\(int\\): 2\n\n"
    testHelperFailure(c, "Assert(1, checker, 2)", nil, true, log,
                      func() interface{} {
        c.Assert(1, checker, 2)
        return nil
    })
}

func (s *HelpersS) TestAssertFailWithExpectedAndMessage(c *gocheck.C) {
    checker := &MyChecker{failCheck: true}
    log := "helpers_test\\.go:[0-9]+ > helpers_test\\.go:[0-9]+:\n" +
           "\\.+ Assert\\(obtained, MyChecker, expected\\):\n" +
           "\\.+ MyObtained \\(int\\): 1\n" +
           "\\.+ MyExpected \\(int\\): 2\n" +
           "\\.+ Hello world!\n\n"
    testHelperFailure(c, "Assert(1, checker, 2, msg)", nil, true, log,
                      func() interface{} {
        c.Assert(1, checker, 2, "Hello", " world!")
        return nil
    })
}

func (s *HelpersS) TestAssertFailWithoutExpected(c *gocheck.C) {
    checker := &MyChecker{failCheck: true, noExpectedValue: true}
    log := "helpers_test\\.go:[0-9]+ > helpers_test\\.go:[0-9]+:\n" +
           "\\.+ Assert\\(obtained, MyChecker\\):\n" +
           "\\.+ MyObtained \\(int\\): 1\n\n"
    testHelperFailure(c, "Assert(1, checker)", nil, true, log,
                      func() interface{} {
        c.Assert(1, checker)
        return nil
    })
}

func (s *HelpersS) TestAssertFailWithoutExpectedAndMessage(c *gocheck.C) {
    checker := &MyChecker{failCheck: true, noExpectedValue: true}
    log := "helpers_test\\.go:[0-9]+ > helpers_test\\.go:[0-9]+:\n" +
           "\\.+ Assert\\(obtained, MyChecker\\):\n" +
           "\\.+ MyObtained \\(int\\): 1\n" +
           "\\.+ Hello world!\n\n"
    testHelperFailure(c, "Assert(1, checker, msg)", nil, true, log,
                      func() interface{} {
        c.Assert(1, checker, "Hello", " world!")
        return nil
    })
}

// -----------------------------------------------------------------------
// Old tests for Assert*() and Check*() helpers not based on checkers.

func (s *HelpersS) TestCheckEqualFailing(c *gocheck.C) {
    log := "helpers_test.go:[0-9]+ > helpers_test.go:[0-9]+:\n" +
           "\\.+ CheckEqual\\(obtained, expected\\):\n" +
           "\\.+ Obtained \\(int\\): 10\n" +
           "\\.+ Expected \\(int\\): 20\n\n"
    testHelperFailure(c, "CheckEqual(10, 20)", false, false, log,
                      func() interface{} {
        return c.CheckEqual(10, 20)
    })
}

func (s *HelpersS) TestCheckEqualFailingWithDiffTypes(c *gocheck.C) {
    log := "helpers_test.go:[0-9]+ > helpers_test.go:[0-9]+:\n" +
           "\\.+ CheckEqual\\(obtained, expected\\):\n" +
           "\\.+ Obtained \\(int\\): 10\n" +
           "\\.+ Expected \\(uint\\): 0xa\n\n"
    testHelperFailure(c, "CheckEqual(10, uint(10))", false, false, log,
                      func() interface{} {
        return c.CheckEqual(10, uint(10))
    })
}

func (s *HelpersS) TestCheckEqualFailingWithNil(c *gocheck.C) {
    log := "helpers_test.go:[0-9]+ > helpers_test.go:[0-9]+:\n" +
           "\\.+ CheckEqual\\(obtained, expected\\):\n" +
           "\\.+ Obtained \\(int\\): 10\n" +
           "\\.+ Expected \\(nil\\): nil\n\n"
    testHelperFailure(c, "CheckEqual(10, nil)", false, false, log,
                      func() interface{} {
        return c.CheckEqual(10, nil)
    })
}

func (s *HelpersS) TestCheckEqualWithMessage(c *gocheck.C) {
    log := "helpers_test.go:[0-9]+ > helpers_test.go:[0-9]+:\n" +
           "\\.+ CheckEqual\\(obtained, expected\\):\n" +
           "\\.+ Obtained \\(int\\): 10\n" +
           "\\.+ Expected \\(int\\): 20\n" +
           "\\.+ That's clearly WRONG!\n\n"
    testHelperFailure(c, "CheckEqual(10, 20, issue)", false, false, log,
                      func() interface{} {
        return c.CheckEqual(10, 20, "That's clearly ", "WRONG!")
    })
}

func (s *HelpersS) TestCheckNotEqualSucceeding(c *gocheck.C) {
    testHelperSuccess(c, "CheckNotEqual(10, 20)", true, func() interface{} {
        return c.CheckNotEqual(10, 20)
    })
}

func (s *HelpersS) TestCheckNotEqualSucceedingWithNil(c *gocheck.C) {
    testHelperSuccess(c, "CheckNotEqual(10, nil)", true, func() interface{} {
        return c.CheckNotEqual(10, nil)
    })
}

func (s *HelpersS) TestCheckNotEqualFailing(c *gocheck.C) {
    log := "helpers_test.go:[0-9]+ > helpers_test.go:[0-9]+:\n" +
           "\\.+ CheckNotEqual\\(obtained, unexpected\\):\n" +
           "\\.+ Both \\(int\\): 10\n\n"
    testHelperFailure(c, "CheckNotEqual(10, 10)", false, false, log,
                      func() interface{} {
        return c.CheckNotEqual(10, 10)
    })
}

func (s *HelpersS) TestCheckNotEqualWithMessage(c *gocheck.C) {
    log := "helpers_test.go:[0-9]+ > helpers_test.go:[0-9]+:\n" +
           "\\.+ CheckNotEqual\\(obtained, unexpected\\):\n" +
           "\\.+ Both \\(int\\): 10\n" +
           "\\.+ That's clearly WRONG!\n\n"
    testHelperFailure(c, "CheckNotEqual(10, 10, issue)", false, false, log,
                      func() interface{} {
        return c.CheckNotEqual(10, 10, "That's clearly ", "WRONG!")
    })
}

func (s *HelpersS) TestAssertEqualSucceeding(c *gocheck.C) {
    testHelperSuccess(c, "AssertEqual(10, 10)", nil, func() interface{} {
        c.AssertEqual(10, 10)
        return nil
    })
}

func (s *HelpersS) TestAssertEqualFailing(c *gocheck.C) {
    log := "helpers_test.go:[0-9]+ > helpers_test.go:[0-9]+:\n" +
           "\\.+ AssertEqual\\(obtained, expected\\):\n" +
           "\\.+ Obtained \\(int\\): 10\n" +
           "\\.+ Expected \\(int\\): 20\n\n"
    testHelperFailure(c, "AssertEqual(10, 20)", nil, true, log,
                      func() interface{} {
        c.AssertEqual(10, 20)
        return nil
    })
}

func (s *HelpersS) TestAssertEqualWithMessage(c *gocheck.C) {
    log := "helpers_test.go:[0-9]+ > helpers_test.go:[0-9]+:\n" +
           "\\.+ AssertEqual\\(obtained, expected\\):\n" +
           "\\.+ Obtained \\(int\\): 10\n" +
           "\\.+ Expected \\(int\\): 20\n" +
           "\\.+ That's clearly WRONG!\n\n"
    testHelperFailure(c, "AssertEqual(10, 20, issue)", nil, true, log,
                      func() interface{} {
        c.AssertEqual(10, 20, "That's clearly ", "WRONG!")
        return nil
    })
}

func (s *HelpersS) TestAssertNotEqualSucceeding(c *gocheck.C) {
    testHelperSuccess(c, "AssertNotEqual(10, 20)", nil, func() interface{} {
        c.AssertNotEqual(10, 20)
        return nil
    })
}

func (s *HelpersS) TestAssertNotEqualFailing(c *gocheck.C) {
    log := "helpers_test.go:[0-9]+ > helpers_test.go:[0-9]+:\n" +
           "\\.+ AssertNotEqual\\(obtained, unexpected\\):\n" +
           "\\.+ Both \\(int\\): 10\n\n"
    testHelperFailure(c, "AssertNotEqual(10, 10)", nil, true, log,
                      func() interface{} {
        c.AssertNotEqual(10, 10)
        return nil
    })
}

func (s *HelpersS) TestAssertNotEqualWithMessage(c *gocheck.C) {
    log := "helpers_test.go:[0-9]+ > helpers_test.go:[0-9]+:\n" +
           "\\.+ AssertNotEqual\\(obtained, unexpected\\):\n" +
           "\\.+ Both \\(int\\): 10\n" +
           "\\.+ That's clearly WRONG!\n\n"
    testHelperFailure(c, "AssertNotEqual(10, 10, issue)", nil, true, log,
                      func() interface{} {
        c.AssertNotEqual(10, 10, "That's clearly ", "WRONG!")
        return nil
    })
}


func (s *HelpersS) TestCheckEqualArraySucceeding(c *gocheck.C) {
    testHelperSuccess(c, "CheckEqual([]byte, []byte)", true, func() interface{} {
        return c.CheckEqual([]byte{1,2}, []byte{1,2})
    })
}

func (s *HelpersS) TestCheckEqualArrayFailing(c *gocheck.C) {
    log := "helpers_test.go:[0-9]+ > helpers_test.go:[0-9]+:\n" +
           "\\.+ CheckEqual\\(obtained, expected\\):\n" +
           "\\.+ Obtained \\(\\[\\]uint8\\): \\[\\]byte{0x1, 0x2}\n" +
           "\\.+ Expected \\(\\[\\]uint8\\): \\[\\]byte{0x1, 0x3}\n\n"
    testHelperFailure(c, "CheckEqual([]byte{2}, []byte{3})", false, false, log,
                      func() interface{} {
        return c.CheckEqual([]byte{1,2}, []byte{1,3})
    })
}

func (s *HelpersS) TestCheckNotEqualArraySucceeding(c *gocheck.C) {
    testHelperSuccess(c, "CheckNotEqual([]byte, []byte)", true,
                      func() interface{} {
        return c.CheckNotEqual([]byte{1,2}, []byte{1,3})
    })
}

func (s *HelpersS) TestCheckNotEqualArrayFailing(c *gocheck.C) {
    log := "helpers_test.go:[0-9]+ > helpers_test.go:[0-9]+:\n" +
           "\\.+ CheckNotEqual\\(obtained, unexpected\\):\n" +
           "\\.+ Both \\(\\[\\]uint8\\): \\[\\]byte{0x1, 0x2}\n\n"
    testHelperFailure(c, "CheckNotEqual([]byte{2}, []byte{3})", false, false,
                      log, func() interface{} {
        return c.CheckNotEqual([]byte{1,2}, []byte{1,2})
    })
}


func (s *HelpersS) TestAssertEqualArraySucceeding(c *gocheck.C) {
    testHelperSuccess(c, "AssertEqual([]byte, []byte)", nil,
                      func() interface{} {
        c.AssertEqual([]byte{1,2}, []byte{1,2})
        return nil
    })
}

func (s *HelpersS) TestAssertEqualArrayFailing(c *gocheck.C) {
    log := "helpers_test.go:[0-9]+ > helpers_test.go:[0-9]+:\n" +
           "\\.+ AssertEqual\\(obtained, expected\\):\n" +
           "\\.+ Obtained \\(\\[\\]uint8\\): \\[\\]byte{0x1, 0x2}\n" +
           "\\.+ Expected \\(\\[\\]uint8\\): \\[\\]byte{0x1, 0x3}\n\n"
    testHelperFailure(c, "AssertEqual([]byte{2}, []byte{3})", nil, true, log,
                      func() interface{} {
        c.AssertEqual([]byte{1,2}, []byte{1,3})
        return nil
    })
}

func (s *HelpersS) TestAssertNotEqualArraySucceeding(c *gocheck.C) {
    testHelperSuccess(c, "AssertNotEqual([]byte, []byte)", nil,
                      func() interface{} {
        c.AssertNotEqual([]byte{1,2}, []byte{1,3})
        return nil
    })
}

func (s *HelpersS) TestAssertNotEqualArrayFailing(c *gocheck.C) {
    log := "helpers_test.go:[0-9]+ > helpers_test.go:[0-9]+:\n" +
           "\\.+ AssertNotEqual\\(obtained, unexpected\\):\n" +
           "\\.+ Both \\(\\[\\]uint8\\): \\[\\]byte{0x1, 0x2}\n\n"
    testHelperFailure(c, "AssertNotEqual([]byte{2}, []byte{3})", nil, true,
                      log, func() interface{} {
        c.AssertNotEqual([]byte{1,2}, []byte{1,2})
        return nil
    })
}

func (s *HelpersS) TestCheckMatchSucceeding(c *gocheck.C) {
    testHelperSuccess(c, "CheckErr('foo', 'fo*')", true, func() interface{} {
        return c.CheckMatch("foo", "fo*")
    })
}

func (s *HelpersS) TestCheckMatchFailing(c *gocheck.C) {
    log := "helpers_test.go:[0-9]+ > helpers_test.go:[0-9]+:\n" +
           "\\.+ CheckMatch\\(value, expression\\):\n" +
           "\\.+ Value \\(string\\): \"foo\"\n" +
           "\\.+ Expected to match expression: \"bar\"\n\n"
    testHelperFailure(c, "CheckMatch('foo', 'bar')", false, false, log,
                      func() interface{} {
        return c.CheckMatch("foo", "bar")
    })
}

func (s *HelpersS) TestCheckMatchFailingWithMessage(c *gocheck.C) {
    log := "helpers_test.go:[0-9]+ > helpers_test.go:[0-9]+:\n" +
           "\\.+ CheckMatch\\(value, expression\\):\n" +
           "\\.+ Value \\(string\\): \"foo\"\n" +
           "\\.+ Expected to match expression: \"bar\"\n" +
           "\\.+ Foo bar!\n\n"
    testHelperFailure(c, "CheckMatch('foo', 'bar')", false, false, log,
                      func() interface{} {
        return c.CheckMatch("foo", "bar", "Foo", " bar!")
    })
}


func (s *HelpersS) TestAssertMatchSucceeding(c *gocheck.C) {
    testHelperSuccess(c, "AssertMatch(s, exp)", nil, func() interface{} {
        c.AssertMatch("str error", "str.*r")
        return nil
    })
}

func (s *HelpersS) TestAssertMatchSucceedingWithError(c *gocheck.C) {
    testHelperSuccess(c, "AssertMatch(os.Error, exp)", nil, func() interface{} {
        c.AssertMatch(os.Errno(13), "perm.*denied")
        return nil
    })
}

func (s *HelpersS) TestAssertMatchFailing(c *gocheck.C) {
    log := "helpers_test.go:[0-9]+ > helpers_test.go:[0-9]+:\n" +
           "\\.+ AssertMatch\\(value, expression\\):\n" +
           "\\.+ Value \\(os\\.Errno\\): 13 \\(permission denied\\)\n" +
           "\\.+ Expected to match expression: \"foo\"\n\n"
    testHelperFailure(c, "AssertMatch(error, foo)", nil, true, log,
                      func() interface{} {
        c.AssertMatch(os.Errno(13), "foo")
        return nil
    })
}

func (s *HelpersS) TestAssertMatchFailingWithPureStrMatch(c *gocheck.C) {
    log := "helpers_test.go:[0-9]+ > helpers_test.go:[0-9]+:\n" +
           "\\.+ AssertMatch\\(value, expression\\):\n" +
           "\\.+ Value \\(string\\): \"foobar\"\n" +
           "\\.+ Expected to match expression: \"foobaz\"\n\n"
    testHelperFailure(c, "AssertMatch('foobar', 'foobaz')", nil, true, log,
                      func() interface{} {
        c.AssertMatch("foobar", "foobaz")
        return nil
    })
}

func (s *HelpersS) TestAssertMatchFailingWithMessage(c *gocheck.C) {
    log := "helpers_test.go:[0-9]+ > helpers_test.go:[0-9]+:\n" +
           "\\.+ AssertMatch\\(value, expression\\):\n" +
           "\\.+ Value \\(string\\): \"foobar\"\n" +
           "\\.+ Expected to match expression: \"foobaz\"\n" +
           "\\.+ That's clearly WRONG!\n\n"
    testHelperFailure(c, "AssertMatch('foobar', 'foobaz')", nil, true, log,
                      func() interface{} {
        c.AssertMatch("foobar", "foobaz", "That's clearly ", "WRONG!")
        return nil
    })
}


// -----------------------------------------------------------------------
// MakeDir() tests.

type MkDirHelper struct {
    path1 string
    path2 string
    isDir1 bool
    isDir2 bool
    isDir3 bool
    isDir4 bool
}

func (s *MkDirHelper) SetUpSuite(c *gocheck.C) {
    s.path1 = c.MkDir()
    s.isDir1 = isDir(s.path1)
}

func (s *MkDirHelper) Test(c *gocheck.C) {
    s.path2 = c.MkDir()
    s.isDir2 = isDir(s.path2)
}

func (s *MkDirHelper) TearDownSuite(c *gocheck.C) {
    s.isDir3 = isDir(s.path1)
    s.isDir4 = isDir(s.path2)
}


func (s *HelpersS) TestMkDir(c *gocheck.C) {
    helper := MkDirHelper{}
    output := String{}
    gocheck.Run(&helper, &gocheck.RunConf{Output: &output})
    c.AssertEqual(output.value, "")
    c.CheckEqual(helper.isDir1, true)
    c.CheckEqual(helper.isDir2, true)
    c.CheckEqual(helper.isDir3, true)
    c.CheckEqual(helper.isDir4, true)
    c.CheckNotEqual(helper.path1, helper.path2)
    c.CheckEqual(isDir(helper.path1), false)
    c.CheckEqual(isDir(helper.path2), false)
}

func isDir(path string) bool {
    if stat, err := os.Stat(path); err == nil {
        return stat.IsDirectory()
    }
    return false
}


// -----------------------------------------------------------------------
// A couple of helper functions to test helper functions. :-)

func testHelperSuccess(c *gocheck.C, name string,
                       expectedResult interface{},
                       closure func() interface{}) {
    var result interface{}
    defer (func() {
        if err := recover(); err != nil {
            panic(err)
        }
        checkState(c, result,
                   &expectedState{
                        name: name,
                        result: expectedResult,
                        failed: false,
                        log: "",
                   })
    })()
    result = closure()
}

func testHelperFailure(c *gocheck.C, name string,
                       expectedResult interface{},
                       shouldStop bool, log string,
                       closure func() interface{}) {
    var result interface{}
    defer (func() {
        if err := recover(); err != nil {
            panic(err)
        }
        checkState(c, result,
                   &expectedState{
                        name: name,
                        result: expectedResult,
                        failed: true,
                        log: log,
                   })
    })()
    result = closure()
    if shouldStop {
        c.Logf("%s didn't stop when it should", name)
    }
}
