package dbee_test

import (
	"dbee/errors"
	"dbee/store"
	"os"

	. "dbee"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var _ = Describe("dbee", func() {
	const defaultPartition string = "default"

	engines := []Engine{Bolt}
	dir := "/tmp/dbee"
	var i store.Instance = nil
	var err error = nil

	Describe("dbee tests", func() {
		for _, eng := range engines {
			Context("Testing engine "+eng.String(), func() {
				_ = os.RemoveAll(dir)
				_ = os.MkdirAll(dir, 0744)

				It("can create instance from scratch", func() {
					i, err = Open(dir, eng)

					Expect(err).NotTo(HaveOccurred())
					Expect(i).ShouldNot(BeNil())
					Expect(i.Dir()).Should(BeIdenticalTo(dir))
				})

				It("will time out if instance is already loaded", func() {
					_, err = Open(dir, eng)

					Expect(err).Should(HaveOccurred())
				})

				It("will close the instance", func() {
					err = i.Close()

					Expect(err).ShouldNot(HaveOccurred())
				})

				It("can re-open the closed instance", func() {
					i, err = Open(dir, eng)

					Expect(err).NotTo(HaveOccurred())
					Expect(i).ShouldNot(BeNil())
					Expect(i.Dir()).Should(BeIdenticalTo(dir))
				})
			})

			Context("Accessing instance propeties/column name", func() {
				var ix uint64
				property := "property"

				It("will retrieve property index", func() {
					ix, err = i.GetPropIndex(property)

					Expect(err).NotTo(HaveOccurred())
					Expect(ix).Should(BeNumerically(">", uint64(0)))
				})

				It("should not find the none existing index", func() {
					_, err := i.GetPropName(0)

					Expect(err).Should(BeIdenticalTo(errors.ErrPropNotFound))
				})

				It("should retrive property name", func() {
					n, err := i.GetPropName(ix)

					Expect(err).NotTo(HaveOccurred())
					Expect(n).Should(BeIdenticalTo(property))
				})

				It("should retrive property index", func() {
					n, err := i.GetPropIndex(property)

					Expect(err).NotTo(HaveOccurred())
					Expect(n).Should(BeIdenticalTo(ix))
				})
			})

			Context("Creating sets", func() {
				DescribeTable("Table for sets", func(n string) {
					s, err := i.Set(n)

					Expect(err).NotTo(HaveOccurred())
					Expect(s).ShouldNot(BeNil())
					Expect(s.Name()).Should(BeIdenticalTo(n))
					Expect(len(s.Partitions())).Should(BeNumerically(">", 0))
					Expect(s.Partitions()[0]).Should(BeIdenticalTo(defaultPartition))
				},
					Entry("set", "set"),
					Entry("set/2", "set/2"),
					Entry("set/2/3", "set/2/3"),
					Entry("set", "set"),
					Entry("set/2", "set/2"),
					Entry("set/2/3", "set/2/3"),
				)

				It("will safely close the instance", func() {
					By("Also clossing all the sets")
					err = i.Close()
					Expect(err).ShouldNot(HaveOccurred())
				})

				It("will safely re-open instance and sets", func() {

					i, err = Open(dir, eng)

					Expect(err).NotTo(HaveOccurred())
					Expect(i).ShouldNot(BeNil())
					Expect(i.Dir()).Should(BeIdenticalTo(dir))

					n := "set"
					s, err := i.Set(n)
					Expect(err).NotTo(HaveOccurred())
					Expect(s).ShouldNot(BeNil())
					Expect(s.Name()).Should(BeIdenticalTo(n))
					Expect(len(s.Partitions())).Should(BeNumerically(">", 0))
					Expect(s.Partitions()[0]).Should(BeIdenticalTo(defaultPartition))
				})
			})

			Context("creating data to sets", func() {
				It("will create a new set transaction", func() {
					set, _ := i.Set("set")
					setTx, err := set.Get()

					Expect(err).NotTo(HaveOccurred())
					Expect(len(setTx.ID())).Should(BeNumerically("==", 26))
					Expect(setTx.CreatedOn().IsZero()).Should(BeTrue())
					Expect(setTx.LastUpdate().IsZero()).Should(BeTrue())
					Expect(setTx.IsSoftDeleted()).Should(BeFalse())
				})
			})

			Context("store and retrive data", func() {
				var setId string

				It("will store a new set", func() {
					set, _ := i.Set("set")
					setTx, _ := set.Get()

					setId = setTx.ID()

					setTx.Wstring("sample", "sample string value")
					setTx.Wstring("prop2", "sample string 2")

					err = setTx.Commit()
					Expect(err).NotTo(HaveOccurred())

					sample := setTx.Rstring("sample")

					prop2 := setTx.Rstring("prop2")

					Expect(sample).Should(BeEquivalentTo("sample string value"))
					Expect(prop2).Should(BeEquivalentTo("sample string 2"))
				})

				It("will retrive the data on id"+setId, func() {
					set, _ := i.Set("set")
					setTx, _ := set.Get(setId)

					sample := setTx.Rstring("sample")
					prop2 := setTx.Rstring("prop2")

					Expect(sample).Should(BeEquivalentTo("sample string value"))
					Expect(prop2).Should(BeEquivalentTo("sample string 2"))
				})

				It("will test all the data type supported", func() {
					set, _ := i.Set("set")
					setTx, _ := set.Get()

					setId = setTx.ID()

					setTx.Wfloat("float", 1)
					setTx.Wdouble("double", 1)
					setTx.Wint("int", 1)
					setTx.Wsint("sint", 1)
					setTx.Wuint("uint", 1)
					setTx.Wbool("bool", true)
					setTx.Wstring("string", "string")
					setTx.Wbytes("bytes", []byte("byte"))

					err = setTx.Commit()

					Expect(err).NotTo(HaveOccurred())
				})

				It("will read all data type supported", func() {
					set, _ := i.Set("set")
					setTx, _ := set.Get(setId)

					Expect(setTx.Rfloat("float")).Should(BeNumerically("==", 1))
					Expect(setTx.Rdouble("double")).Should(BeNumerically("==", 1))
					Expect(setTx.Rint("int")).Should(BeNumerically("==", 1))
					Expect(setTx.Rsint("sint")).Should(BeNumerically("==", 1))
					Expect(setTx.Ruint("uint")).Should(BeNumerically("==", 1))
					Expect(setTx.Rbool("bool")).Should(BeTrue())
					Expect(setTx.Rstring("string")).Should(BeIdenticalTo("string"))
					Expect(setTx.Rbytes("bytes")).Should(Equal([]byte("byte")))
				})
			})

			Context("data deletion", func() {
				var setId string

				It("will a soft delete", func() {
					set, _ := i.Set("set")
					setTx, _ := set.Get()
					setId = setTx.ID()
					setTx.Wint("int", 1)
					_ = setTx.Commit()

					Expect(setTx.IsSoftDeleted()).Should(BeFalse())

					setTx.Delete()

					err = setTx.Commit()
					Expect(err).NotTo(HaveOccurred())
				})

				It("will retrive the soft deleted entry", func() {
					set, _ := i.Set("set")
					setTx, _ := set.Get(setId)

					Expect(setTx.IsSoftDeleted()).Should(BeTrue())
					Expect(setTx.Rint("int")).Should(BeNumerically("==", 1))
				})

				It("will hard delete the data", func() {
					set, _ := i.Set("set")
					setTx, _ := set.Get(setId)
					err = setTx.HardDelete()

					Expect(err).NotTo(HaveOccurred())
				})

				It("will check if the data is hardly deleted", func() {
					set, _ := i.Set("set")
					setTx, _ := set.Get(setId)

					Expect(setTx.Rint("int")).Should(BeNumerically("==", 0))
				})
			})

			Context("partition access", func() {
				It("willl verify the default partition", func() {
					set, _ := i.Set("set")
					setTx, _ := set.Get()

					Expect(setTx.Partition().Name()).
						Should(BeEquivalentTo("default"))
				})

				It("will create a new partition", func() {
					set, _ := i.Set("set")
					par, err := set.Partition("new partition")

					Expect(err).NotTo(HaveOccurred())
					Expect(par.Name()).Should(BeEquivalentTo("new partition"))
				})

				It("will create a new tx from the new partition", func() {
					set, _ := i.Set("set")
					par, _ := set.Partition("new partition")
					setTx, err := par.Get()

					Expect(err).NotTo(HaveOccurred())
					Expect(len(setTx.ID())).Should(BeNumerically("==", 26))
				})
			})

			Context("store and retrive data in the new partition", func() {
				var setId string
				pname := "new partition"

				It("will store a new set", func() {
					set, _ := i.Set("set")
					p, _ := set.Partition(pname)
					setTx, _ := p.Get()

					setId = setTx.ID()

					setTx.Wstring("sample", "sample string value")
					setTx.Wstring("prop2", "sample string 2")

					err = setTx.Commit()
					Expect(err).NotTo(HaveOccurred())

					sample := setTx.Rstring("sample")

					prop2 := setTx.Rstring("prop2")

					Expect(sample).Should(BeEquivalentTo("sample string value"))
					Expect(prop2).Should(BeEquivalentTo("sample string 2"))
				})

				It("will retrive the data on id"+setId, func() {
					set, _ := i.Set("set")
					p, _ := set.Partition(pname)
					setTx, _ := p.Get(setId)

					sample := setTx.Rstring("sample")
					prop2 := setTx.Rstring("prop2")

					Expect(sample).Should(BeEquivalentTo("sample string value"))
					Expect(prop2).Should(BeEquivalentTo("sample string 2"))
				})

				It("will test all the data type supported", func() {
					set, _ := i.Set("set")
					p, _ := set.Partition(pname)
					setTx, _ := p.Get()

					setId = setTx.ID()

					setTx.Wfloat("float", 1)
					setTx.Wdouble("double", 1)
					setTx.Wint("int", 1)
					setTx.Wsint("sint", 1)
					setTx.Wuint("uint", 1)
					setTx.Wbool("bool", true)
					setTx.Wstring("string", "string")
					setTx.Wbytes("bytes", []byte("byte"))

					err = setTx.Commit()

					Expect(err).NotTo(HaveOccurred())
				})

				It("will read all data type supported", func() {
					set, _ := i.Set("set")
					p, _ := set.Partition(pname)
					setTx, _ := p.Get(setId)

					Expect(setTx.Rfloat("float")).Should(BeNumerically("==", 1))
					Expect(setTx.Rdouble("double")).Should(BeNumerically("==", 1))
					Expect(setTx.Rint("int")).Should(BeNumerically("==", 1))
					Expect(setTx.Rsint("sint")).Should(BeNumerically("==", 1))
					Expect(setTx.Ruint("uint")).Should(BeNumerically("==", 1))
					Expect(setTx.Rbool("bool")).Should(BeTrue())
					Expect(setTx.Rstring("string")).Should(BeIdenticalTo("string"))
					Expect(setTx.Rbytes("bytes")).Should(Equal([]byte("byte")))
				})
			})

			Context("data deletion the new partition", func() {
				var setId string

				It("will a soft delete", func() {
					set, _ := i.Set("set")
					p, _ := set.Partition("new partition")
					setTx, _ := p.Get()

					setId = setTx.ID()
					setTx.Wint("int", 1)
					_ = setTx.Commit()

					Expect(setTx.IsSoftDeleted()).Should(BeFalse())

					setTx.Delete()

					err = setTx.Commit()
					Expect(err).NotTo(HaveOccurred())
				})

				It("will retrive the soft deleted entry", func() {
					set, _ := i.Set("set")
					p, _ := set.Partition("new partition")
					setTx, _ := p.Get(setId)

					Expect(setTx.IsSoftDeleted()).Should(BeTrue())
					Expect(setTx.Rint("int")).Should(BeNumerically("==", 1))
				})

				It("will hard delete the data", func() {
					set, _ := i.Set("set")
					p, _ := set.Partition("new partition")
					setTx, _ := p.Get(setId)
					err = setTx.HardDelete()

					Expect(err).NotTo(HaveOccurred())
				})

				It("will check if the data is hardly deleted", func() {
					set, _ := i.Set("set")
					p, _ := set.Partition("new partition")
					setTx, _ := p.Get(setId)

					Expect(setTx.Rint("int")).Should(BeNumerically("==", 0))
				})
			})

			Context("Check indexing", func() {
				idx := []string{
					"index.first",
					"index.second",
					"index.third",
					"username",
				}

				It("will check that the indexing is empty at first", func() {
					set, _ := i.Set("set")
					l := len(set.ListIndexes())

					Expect(l).Should(BeNumerically("==", 0))
				})

				It("will setup a new index", func() {
					set, _ := i.Set("set")

					err = set.Index(idx[0])

					Expect(err).NotTo(HaveOccurred())

					_ = set.Index(idx[1])
					_ = set.Index(idx[2])
					_ = set.Index(idx[3])

					l := len(set.ListIndexes())
					Expect(l).Should(BeNumerically("==", 4))
					Expect(set.ListIndexes()).Should(ConsistOf(idx))
				})

				It("will close the db and recheck the indexes", func() {
					err = i.Close()
					Expect(err).NotTo(HaveOccurred())

					i, err = Open(dir, eng)
					Expect(err).NotTo(HaveOccurred())

					set, err := i.Set("set")
					Expect(err).NotTo(HaveOccurred())

					l := len(set.ListIndexes())
					Expect(l).Should(BeNumerically("==", 4))
					Expect(set.ListIndexes()).Should(ConsistOf(idx))
				})
			})
		}
	})
})
