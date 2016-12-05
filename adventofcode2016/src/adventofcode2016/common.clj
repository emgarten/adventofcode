(ns adventofcode2016.common
  (:use [clojure.pprint])
  (:require [clojure.java.io :as io])
  (:require [clojure.string :as str]))

(defn remove-empty
  [items]
  (filter #(> (count %) 0)
    (map #(str/trim %) items)))

(defn get-lines
  "Parse and trim lines of text."
  [input-data]
  (remove-empty
   (str/split-lines input-data)))

(defn slurp-resource
  "Read a file and split by lines"
  [file-name]
  (-> file-name
      io/resource
      io/file
      slurp))
