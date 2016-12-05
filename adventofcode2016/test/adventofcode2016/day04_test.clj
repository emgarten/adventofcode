(ns adventofcode2016.day04-test
  (:use [adventofcode2016.common])
  (:require [clojure.test :refer :all]
            [adventofcode2016.day04 :refer :all]))

(deftest solve-part-1-test
  (testing "first example"
    (is (= 123 (solve-part-1 "aaaaa-bbb-z-y-x-123[abxyz]"))))
  (testing "first example modified"
    (is (= 123 (solve-part-1 "xyz-aaaaa-bbb-123[abxyz]"))))
  (testing "non-alpha"
    (is (= 404 (solve-part-1 "not-a-real-room-404[oarel]"))))
  (testing "false"
    (is (= 0 (solve-part-1 "totally-real-room-200[decoy]"))))
  (testing "above max"
    (is (= 987 (solve-part-1 "a-b-c-d-e-f-g-h-987[abcde]"))))
  (testing "above max with extra count"
    (is (= 987 (solve-part-1 "a-b-c-d-e-f-g-hh-987[habcd]"))))
  (testing "multi-line"
    (is (= 1110 (solve-part-1 "aaaaa-bbb-z-y-x-123[abxyz]\na-b-c-d-e-f-g-h-987[abcde]\n"))))
  (testing "ordering"
    (is (= 123 (solve-part-1 "yxffaaaddddd-123[dafxy]"))))
  (testing "ordering ties"
    (is (= 123 (solve-part-1 "tt-bb-zzz-yyy-xxx-aaa-sss-123[asxyz]")))))
