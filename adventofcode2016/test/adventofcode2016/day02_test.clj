(ns adventofcode2016.day02-test
  (:require [clojure.test :refer :all]
            [adventofcode2016.day02 :refer :all]))

(deftest solve-part-1-test
  (testing "Single line"
    (is (= '(1) (solve-part-1 "ULL"))))
  (testing "Four lines"
    (is (= '(1 9 8 5) (solve-part-1 "ULL
                               RRDDD
                               LURDL
                               UUUUD \n "))))
  (testing "Complex"
    (is (= '(9 5) (solve-part-1 "DRRRRRRRRRRRRRUUUUUUUDUDUDUDUDULLDDRR\nULLR")))))

(deftest process-single-key-test
  (testing "9 -> L -> 8"
    (is (= 21 (process-single-key "L" 22))))
  (testing "8 -> U -> 5"
    (is (= 11 (process-single-key "U" 21))))
  (testing "5 -> R -> 6"
    (is (= 12 (process-single-key "R" 11))))
  (testing "6 -> D -> 9"
    (is (= 22 (process-single-key "D" 12))))
  (testing "9 -> L -> 8"
    (is (= 21 (process-single-key "L" 22)))))

(deftest from-key-id-test
  (testing "1"
    (is (= 1 (from-key-id 0))))
  (testing "2"
    (is (= 2 (from-key-id 1))))
  (testing "3"
    (is (= 3 (from-key-id 2))))
  (testing "4"
    (is (= 4 (from-key-id 10))))
  (testing "5"
    (is (= 5 (from-key-id 11))))
  (testing "6"
    (is (= 6 (from-key-id 12)))))
