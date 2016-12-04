(ns adventofcode2016.day03
  (:use [adventofcode2016.common])
  (:use [clojure.pprint])
  (:require [clojure.string :as str]))

(defn parse-line
  [input-data]
  (map #(Integer. %)
       (filter #(> (count %) 0)
         (map #(str/trim %)
              (str/split input-data #" ")))))

(defn get-triples
  [lines]
  (map #(parse-line %) lines))

(defn test-side
  [total side]
  (let [result (> (- total side) side)]
    result))

(defn triangle?
  [lengths]
  (let [total (reduce + lengths)]
    (every? true? (map #(test-side total %) lengths))))

(defn solve-part-1
  [input-data]
  (count
   (filter #(triangle? %)
     (get-triples
      (get-lines input-data)))))
