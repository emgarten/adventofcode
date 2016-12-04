(ns adventofcode2016.day03-test
  (:require [clojure.test :refer :all]
            [adventofcode2016.day03 :refer :all]))

(deftest solve-part-1-test
  (testing "Single line false"
    (is (= 0 (solve-part-1 "  5  10  25  "))))
  (testing "Single line true"
    (is (= 1 (solve-part-1 "  1 1 1  "))))
  (testing "Single line true also"
    (is (= 1 (solve-part-1 "  1 10 10  ")))))
