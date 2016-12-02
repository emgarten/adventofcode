(ns adventofcode2016.core
  (:require [adventofcode2016.day01 :as day01])
  (:require [clojure.java.io :as io])
  (:use [clojure.pprint]))

(defn main
  []
  (-> "day01.txt"
      io/resource
      io/file
      slurp
      day01/solve-part-2
      pprint))
