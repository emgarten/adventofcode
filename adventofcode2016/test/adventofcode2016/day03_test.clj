(ns adventofcode2016.day03-test
  (:use [adventofcode2016.common])
  (:require [clojure.test :refer :all]
            [adventofcode2016.day03 :refer :all]))

(deftest solve-part-1-test
  (testing "Single line false"
    (is (= 0 (solve-part-1 "  5  10  25  "))))
  (testing "Single line true"
    (is (= 1 (solve-part-1 "  1 1 1  "))))
  (testing "Single line true also"
    (is (= 1 (solve-part-1 "  1 10 10  ")))))

(deftest get-vertical-triples-test
  (testing "basic"
    (is (= '(101 102 103) (first (get-vertical-triples (get-lines "101 301 501
                            102 302 502
                            103 303 503
                            201 401 601
                            202 402 602
                            203 403 603 ")))))))
