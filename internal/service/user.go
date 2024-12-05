package service

//func (s Service) GetUserByID(ctx context.Context, id string) (*contract.UserResponse, error) {
//	user, err := s.storage.GetUserByID(ctx, id)
//	if err != nil {
//		if errors.Is(err, db.ErrNotFound) {
//			return nil, contract.ErrInvalidRequest
//		}
//		return nil, fmt.Errorf("get user: %w", err)
//	}
//
//	return &contract.UserResponse{
//		Id:        user.Id,
//		Name:      user.Name,
//		Email:     user.Email,
//		CreatedAt: user.CreatedAt,
//	}, nil
//}
//
//func (s Service) CreateUser(ctx context.Context, req *contract.UserRequest) (*contract.UserResponse, error) {
//	newUser := db.User{
//		Id:        nanoid.Must(8),
//		Name:      req.Name,
//		Email:     req.Email,
//		CreatedAt: req.CreatedAt,
//	}
//
//	if err := s.storage.CreateUser(ctx, newUser); err != nil {
//		return nil, fmt.Errorf("create user: %w", err)
//	}
//
//	res, err := s.storage.GetUserByID(ctx, newUser.Id)
//
//	if err != nil {
//		return nil, fmt.Errorf("get user: %w", err)
//	}
//
//	return &contract.UserResponse{
//		Id:        res.Id,
//		Name:      res.Name,
//		Email:     res.Email,
//		CreatedAt: res.CreatedAt,
//	}, nil
//}
//
//func (s Service) UpdateUser(ctx context.Context, req *contract.UserRequest) error {
//	user := db.User{
//		Id: req.Id,
//	}
//
//	err := s.storage.UpdateUser(ctx, user)
//	if err != nil {
//		return fmt.Errorf("update user: %w", err)
//	}
//
//	return nil
//}
//
//func (s Service) DeleteUser(ctx context.Context, id string) error {
//	err := s.storage.DeleteUser(ctx, id)
//	if err != nil {
//		return fmt.Errorf("delete user: %w", err)
//	}
//
//	return nil
//}
//
//func (s Service) GetAllUsers(ctx context.Context) ([]contract.UserResponse, error) {
//	users, err := s.storage.GetAllUsers(ctx)
//	if err != nil {
//		return nil, fmt.Errorf("get all users: %w", err)
//	}
//
//	usersResponse := make([]contract.UserResponse, 0)
//
//	for _, user := range users {
//		userResponse := contract.UserResponse{
//			Id:        user.Id,
//			Name:      user.Name,
//			Email:     user.Email,
//			CreatedAt: user.CreatedAt,
//		}
//		usersResponse = append(usersResponse, userResponse)
//	}
//
//	return usersResponse, nil
//}
