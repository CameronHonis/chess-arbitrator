package matchmaking_test

import (
	"github.com/CameronHonis/chess-arbitrator/builders"
	"github.com/CameronHonis/chess-arbitrator/helpers/mocks"
	"github.com/CameronHonis/chess-arbitrator/matchmaking"
	"github.com/CameronHonis/chess-arbitrator/models"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.uber.org/mock/gomock"
)

func CreateServices(ctrl *gomock.Controller) *matchmaking.MatchmakingService {
	logServiceMock := mocks.NewMockLoggerServiceI(ctrl)
	logServiceMock.EXPECT().Log(gomock.Any(), gomock.Any()).AnyTimes()
	logServiceMock.EXPECT().SetParent(gomock.Any()).AnyTimes()

	matchServiceMock := mocks.NewMockMatcherServiceI(ctrl)
	matchServiceMock.EXPECT().SetParent(gomock.Any()).AnyTimes()

	matchmakingService := matchmaking.NewMatchmakingService(matchmaking.NewMatchmakingConfig())
	matchmakingService.AddDependency(logServiceMock)
	matchmakingService.AddDependency(matchServiceMock)
	return matchmakingService
}

var _ = Describe("MatchmakingService", func() {
	var matchmakingService *matchmaking.MatchmakingService
	BeforeEach(func() {
		ctrl := gomock.NewController(T, gomock.WithOverridableExpectations())
		matchmakingService = CreateServices(ctrl)
	})
	Describe("::AddClient", func() {
		var client *models.ClientProfile
		var timeControl *models.TimeControl
		BeforeEach(func() {
			client = models.NewClientProfile("some-client-key", 1000)
			timeControl = builders.NewBlitzTimeControl()
		})
		When("the client is not already in the pool", func() {
			It("should add the client to the pool dedicated to the time control", func() {
				Expect(matchmakingService.AddClient(client, timeControl)).To(Succeed())
				Expect(matchmakingService.GetClientCountByTimeControl(timeControl)).To(Equal(1))
				Expect(matchmakingService.GetClientCountByTimeControl(builders.NewBulletTimeControl())).To(Equal(0))
			})
		})
		When("the client already exists in a pool", func() {
			BeforeEach(func() {
				Expect(matchmakingService.AddClient(client, timeControl)).To(Succeed())
			})
			Context("and the time control is the same as before", func() {
				It("should return an error", func() {
					Expect(matchmakingService.AddClient(client, timeControl)).To(HaveOccurred())
				})
			})
			Context("and the time control is different than before", func() {
				It("should return an error", func() {
					Expect(matchmakingService.AddClient(client, builders.NewBulletTimeControl())).To(HaveOccurred())
				})
			})
		})
	})
	Describe("::RemoveClient", func() {
		var client *models.ClientProfile
		var timeControl *models.TimeControl
		BeforeEach(func() {
			client = models.NewClientProfile("some-client-key", 1000)
			timeControl = builders.NewBlitzTimeControl()
		})
		When("the client is in the pool", func() {
			BeforeEach(func() {
				Expect(matchmakingService.AddClient(client, timeControl)).To(Succeed())
				Expect(matchmakingService.GetClientCountByTimeControl(timeControl)).To(Equal(1))
			})
			It("removes the client from the pool", func() {
				Expect(matchmakingService.RemoveClient(client.ClientKey)).To(Succeed())
				Expect(matchmakingService.GetClientCountByTimeControl(timeControl)).To(Equal(0))
			})
		})
		When("the client is not in the pool", func() {
			It("returns an error", func() {
				Expect(matchmakingService.RemoveClient(client.ClientKey)).To(HaveOccurred())
			})
		})
	})
})
