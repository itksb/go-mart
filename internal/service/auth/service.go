package auth

import (
	"context"
	"github.com/itksb/go-mart/internal/service/auth/token"
	"github.com/itksb/go-mart/internal/service/auth/validation"
	"github.com/itksb/go-mart/pkg/logger"
	"time"
)

type Service struct {
	identity     IdentityInterface
	log          logger.Interface
	hashAlgo     HashAlgoInterface
	tokenCreate  token.CreateTokenFunc
	tokenParse   token.ParseWithClaimsFunc
	secretReader token.Secret
	now          func() time.Time
}

type Opts struct {
	IdentityProvider IdentityInterface
	Logger           logger.Interface
	HashAlgo         HashAlgoInterface
	TokenCreate      token.CreateTokenFunc
	TokenParse       token.ParseWithClaimsFunc
	SecretReader     token.Secret
	NowTime          func() time.Time
}

func NewAuthService(opts Opts) (*Service, error) {
	if opts.IdentityProvider == nil {
		return nil, ErrOptsIdentityProviderIsNil
	}
	if opts.Logger == nil {
		return nil, ErrOptsLoggerIsNil
	}
	if opts.HashAlgo == nil {
		return nil, ErrOptsLoggerIsNil
	}
	if opts.TokenCreate == nil {
		return nil, ErrOptsTokenCreateIsNil
	}
	if opts.TokenCreate == nil {
		return nil, ErrOptsTokenCreateIsNil
	}
	if opts.TokenParse == nil {
		return nil, ErrOptsTokenParseIsNil
	}
	if opts.SecretReader == nil {
		return nil, ErrOptsSecretIsNil
	}
	if opts.NowTime == nil {
		return nil, ErrOptsNowTimeIsNil
	}

	return &Service{
		identity:    opts.IdentityProvider,
		log:         opts.Logger,
		hashAlgo:    opts.HashAlgo,
		tokenCreate: opts.TokenCreate,
		tokenParse:  opts.TokenParse,
		now:         opts.NowTime,
	}, nil
}

const HashAlgoDefaultCost = 10

func (s *Service) SignUp(ctx context.Context, cred ClientCredential) (newToken string, err error) {
	// it is nothing criminal to validate input data twice: in controller layer and here in service
	if _, err := validation.ValidatePassword(cred.Password); err != nil {
		s.log.Infof("password validation fails. Err: %s", err.Error())
		return "", err
	}
	if _, err := validation.ValidateLogin(cred.Login); err != nil {
		s.log.Infof("login validation fails. Err: %s", err.Error())
		return "", err
	}

	hash, err := s.hashAlgo.GenerateFromPassword([]byte(cred.Password), HashAlgoDefaultCost)
	if err != nil {
		s.log.Errorf("hash generation error. Err: %s", err.Error())
		return "", err
	}

	identUser, err := s.identity.Create(ctx, cred.Login, string(hash))
	if err != nil {
		s.log.Errorf("identity creation error. Err: %s", err.Error())
		return "", err
	}

	newToken, err = s.tokenCreate(
		&token.MartClaims{
			UserID:    identUser.ID,
			ExpiresAt: s.now().Add(24 * time.Hour),
		}, s.secretReader)
	if err != nil {
		s.log.Errorf("token creating fails. Err: %s", err.Error())
		return "", err
	}

	return newToken, nil
}

func (s *Service) SignIn(ctx context.Context, cred ClientCredential) (newToken string, err error) {
	// it is nothing criminal to validate input data twice: in controller layer and here in service
	if _, err := validation.ValidatePassword(cred.Password); err != nil {
		s.log.Infof("password validation fails. Err: %s", err.Error())
		return "", err
	}
	if _, err := validation.ValidateLogin(cred.Login); err != nil {
		s.log.Infof("login validation fails. Err: %s", err.Error())
		return "", err
	}

	identUser, err := s.identity.FindOne(ctx, cred.Login, "")
	if err != nil {
		s.log.Infof("login %s not found. Error: %s", cred.Login, err.Error())
		return "", ErrInvalidCredentials
	}

	err = s.hashAlgo.CompareHashAndPassword([]byte(identUser.PasswordHash), []byte(cred.Password))
	if err != nil {
		s.log.Infof("password for login %s is wrong. Err: %s", cred.Login, err.Error())
		return "", ErrInvalidCredentials
	}

	newToken, err = s.tokenCreate(
		&token.MartClaims{
			UserID:    identUser.ID,
			ExpiresAt: s.now().Add(24 * time.Hour),
		}, s.secretReader)

	if err != nil {
		s.log.Errorf("token creating fails. Err: %s", err.Error())
		return "", err
	}

	return newToken, nil
}
