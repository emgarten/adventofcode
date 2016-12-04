(ns adventofcode2016.core
  (:require [adventofcode2016.day03 :as day03])
  (:require [clojure.java.io :as io])
  (:use [clojure.pprint]))

(defn main
  []
  (-> "day03.txt"
      io/resource
      io/file
      slurp
      day03/solve-part-2
      pprint))
