(ns adventofcode2016.common
  (:use [clojure.pprint])
  (:require [clojure.string :as str]))

(defn get-lines
  "Parse and trim lines of text."
  [input-data]
  (filter #(> (count %) 0)
    (map #(str/trim %)
      (str/split-lines input-data))))
