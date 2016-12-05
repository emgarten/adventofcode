(ns adventofcode2016.core
  (:require [adventofcode2016.day04 :as day04])
  (:use [adventofcode2016.common])
  (:require [clojure.java.io :as io])
  (:use [clojure.pprint]))

(defn main
  []
  (-> "day04.txt"
      slurp-resource
      day04/solve-part-1
      pprint))
